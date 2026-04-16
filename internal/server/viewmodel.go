package server

import (
	"goserviceman/internal/config"
	"goserviceman/internal/process"
	"goserviceman/internal/views"
)

// BuildPageVM constructs the view model for the full dashboard.
func BuildPageVM(apps []config.AppConfig, statuses map[string]process.Status) views.PageVM {
	vms := make([]views.AppVM, 0, len(apps))
	for _, a := range apps {
		vms = append(vms, BuildAppVM(a, statuses[a.Name]))
	}
	return views.PageVM{Apps: vms}
}

// BuildAppVM constructs the view model for a single app row.
func BuildAppVM(app config.AppConfig, status process.Status) views.AppVM {
	return views.AppVM{
		Name:       app.Name,
		Executable: app.Executable,
		LogFile:    app.LogFile,
		Port:       app.Port,
		Running:    status == process.StatusRunning,
		StatusText: status.String(),
	}
}
