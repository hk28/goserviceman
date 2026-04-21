package server

import (
	"embed"
	"fmt"
	"net/http"

	"goserviceman/internal/config"
	"goserviceman/internal/process"

	"github.com/a-h/templ"
)

// Server holds the shared dependencies for all HTTP handlers.
type Server struct {
	mgr         *process.ProcManager
	apps        []config.AppConfig
	cfgPath     string
	staticFiles embed.FS
}

// New creates a Server.
func New(mgr *process.ProcManager, apps []config.AppConfig, cfgPath string, staticFiles embed.FS) *Server {
	return &Server{mgr: mgr, apps: apps, cfgPath: cfgPath, staticFiles: staticFiles}
}

// Handler builds and returns the HTTP router.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.handleIndex)
	mux.Handle("GET /static/", http.FileServer(http.FS(s.staticFiles)))

	mux.HandleFunc("GET /app/{name}/status", s.handleStatus)
	mux.HandleFunc("GET /app/{name}/logs", s.handleLogs)
	mux.HandleFunc("POST /app/{name}/start", s.handleStart)
	mux.HandleFunc("POST /app/{name}/stop", s.handleStop)
	mux.HandleFunc("POST /app/{name}/restart", s.handleRestart)
	mux.HandleFunc("POST /app/{name}/open-browser", s.handleOpenBrowser)

	mux.HandleFunc("POST /reload-config", s.handleReloadConfig)

	return mux
}

// render writes a templ component to w.
func render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `<div style="color:red">render error: %v</div>`, err)
	}
}
