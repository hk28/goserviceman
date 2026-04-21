package process

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"goserviceman/internal/config"

	"github.com/mitchellh/go-ps"
)

// Status represents whether a managed process is running or stopped.
type Status int

const (
	StatusStopped Status = iota
	StatusRunning
)

func (s Status) String() string {
	if s == StatusRunning {
		return "running"
	}
	return "stopped"
}

type procEntry struct {
	PID int
	Cmd *exec.Cmd // nil when reattached via go-ps (process not owned by us)
}

// ProcManager manages the lifecycle of configured applications.
type ProcManager struct {
	mu      sync.Mutex
	entries map[string]*procEntry // key = AppConfig.Name
	apps    map[string]config.AppConfig
}

// New creates a new ProcManager for the given app configs.
func New(apps []config.AppConfig) *ProcManager {
	m := &ProcManager{
		entries: make(map[string]*procEntry),
		apps:    make(map[string]config.AppConfig),
	}
	for _, a := range apps {
		m.apps[a.Name] = a
	}
	return m
}

// Reattach scans running processes and reattaches to any that match configured apps.
// Call once at startup so pre-existing processes show as running.
func (m *ProcManager) Reattach() {
	procs, err := ps.Processes()
	if err != nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, app := range m.apps {
		want := exeBase(app.Executable)
		for _, p := range procs {
			if exeBase(p.Executable()) == want {
				m.entries[app.Name] = &procEntry{PID: p.Pid()}
				break
			}
		}
	}
}

// StartImmediately starts any configured apps marked for automatic startup.
func (m *ProcManager) StartImmediately() {
	for _, app := range m.apps {
		if !app.StartImmediately {
			continue
		}
		if m.StatusOf(app.Name) == StatusRunning {
			continue
		}
		if err := m.Start(app.Name); err != nil {
			log.Printf("auto-start %s: %v", app.Name, err)
		} else {
			log.Printf("auto-started %s", app.Name)
		}
	}
}

// Start launches the named application.
func (m *ProcManager) Start(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	app, ok := m.apps[name]
	if !ok {
		return fmt.Errorf("unknown app: %s", name)
	}
	if e, running := m.entries[name]; running && isAlive(e.PID) {
		return fmt.Errorf("%s is already running (pid %d)", name, e.PID)
	}

	cmd := exec.Command(app.Executable, app.Args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting %s: %w", name, err)
	}
	log.Printf("started %s (pid %d)", name, cmd.Process.Pid)

	entry := &procEntry{PID: cmd.Process.Pid, Cmd: cmd}
	m.entries[name] = entry

	go func() {
		cmd.Wait() //nolint:errcheck
		m.mu.Lock()
		// Only clear if this is still the same entry (not replaced by a restart)
		if cur, ok := m.entries[name]; ok && cur == entry {
			delete(m.entries, name)
		}
		m.mu.Unlock()
	}()

	return nil
}

// Stop terminates the named application.
func (m *ProcManager) Stop(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry, ok := m.entries[name]
	if !ok {
		return fmt.Errorf("%s is not running", name)
	}
	if err := killProc(entry); err != nil {
		return fmt.Errorf("stopping %s: %w", name, err)
	}
	delete(m.entries, name)
	return nil
}

// Restart stops then starts the named application.
func (m *ProcManager) Restart(name string) error {
	// Ignore stop errors (process may already be dead)
	m.Stop(name) //nolint:errcheck
	return m.Start(name)
}

// StatusOf returns the current status of the named application.
// It also reconciles stale entries by checking liveness via go-ps.
func (m *ProcManager) StatusOf(name string) Status {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry, ok := m.entries[name]
	if !ok {
		return StatusStopped
	}
	if !isAlive(entry.PID) {
		delete(m.entries, name)
		return StatusStopped
	}
	return StatusRunning
}

// AllStatuses returns the current status of every configured application.
func (m *ProcManager) AllStatuses() map[string]Status {
	out := make(map[string]Status, len(m.apps))
	for name := range m.apps {
		out[name] = m.StatusOf(name)
	}
	return out
}

// isAlive returns true if a process with the given PID is still running.
func isAlive(pid int) bool {
	p, err := ps.FindProcess(pid)
	return err == nil && p != nil
}

// exeBase returns the base filename of an executable path, without extension.
func exeBase(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, ".exe")
}
