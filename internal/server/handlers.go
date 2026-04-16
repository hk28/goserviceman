package server

import (
	"fmt"
	"log"
	"net/http"

	"goserviceman/internal/browser"
	"goserviceman/internal/config"
	"goserviceman/internal/logs"
	"goserviceman/internal/process"
	"goserviceman/internal/views"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	vm := BuildPageVM(s.apps, s.mgr.AllStatuses())
	render(w, r, views.IndexPage(vm))
}

func (s *Server) handleLogs(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	app, ok := s.appByName(name)
	if !ok {
		http.NotFound(w, r)
		return
	}
	vm := views.LogVM{AppName: name}
	if app.LogFile != "" {
		lines, err := logs.TailN(app.LogFile, 100)
		if err != nil {
			vm.Error = err.Error()
		} else {
			vm.Lines = lines
		}
	} else {
		vm.Error = "no log file configured for this app"
	}
	render(w, r, views.LogPanel(vm))
}

func (s *Server) handleStart(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if err := s.mgr.Start(name); err != nil {
		log.Printf("start %s: %v", name, err)
	}
	s.renderAppRow(w, r, name)
}

func (s *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if err := s.mgr.Stop(name); err != nil {
		log.Printf("stop %s: %v", name, err)
	}
	s.renderAppRow(w, r, name)
}

func (s *Server) handleRestart(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if err := s.mgr.Restart(name); err != nil {
		log.Printf("restart %s: %v", name, err)
	}
	s.renderAppRow(w, r, name)
}

func (s *Server) handleOpenBrowser(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	app, ok := s.appByName(name)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if app.Port <= 0 {
		fmt.Fprint(w, `<span class="text-[#f87171] text-[11px]">no port configured</span>`)
		return
	}
	url := fmt.Sprintf("http://localhost:%d", app.Port)
	if err := browser.Open(url); err != nil {
		log.Printf("open browser for %s: %v", name, err)
		fmt.Fprintf(w, `<span class="text-[#f87171] text-[11px]">failed: %v</span>`, err)
		return
	}
	fmt.Fprintf(w, `<span class="text-[#34d399] text-[11px]">opened ✓</span>`)
}

func (s *Server) handleReloadConfig(w http.ResponseWriter, r *http.Request) {
	apps, err := config.Load(s.cfgPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.apps = apps
	s.mgr = process.New(apps)
	s.mgr.Reattach()
	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusNoContent)
}

// renderAppRow re-renders the row for one app after an action.
func (s *Server) renderAppRow(w http.ResponseWriter, r *http.Request, name string) {
	app, ok := s.appByName(name)
	if !ok {
		http.NotFound(w, r)
		return
	}
	vm := BuildAppVM(app, s.mgr.StatusOf(name))
	render(w, r, views.AppRow(vm))
}

// appByName finds an AppConfig by name.
func (s *Server) appByName(name string) (config.AppConfig, bool) {
	for _, a := range s.apps {
		if a.Name == name {
			return a, true
		}
	}
	return config.AppConfig{}, false
}
