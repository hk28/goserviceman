package main

//go:generate templ generate

import (
	"embed"
	"flag"
	"log"
	"net/http"

	"goserviceman/internal/config"
	"goserviceman/internal/process"
	"goserviceman/internal/server"
)

//go:embed static
var staticFiles embed.FS

func main() {
	cfgPath := flag.String("config", "goserviceman.yaml", "path to config YAML")
	addr := flag.String("addr", ":9000", "listen address")
	flag.Parse()

	apps, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	mgr := process.New(apps)
	mgr.Reattach()

	srv := server.New(mgr, apps, *cfgPath, staticFiles)

	log.Printf("goserviceman listening on http://localhost%s", *addr)
	if err := http.ListenAndServe(*addr, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}
