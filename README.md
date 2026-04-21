# goserviceman

A lightweight web-based process manager. Configure your apps in a YAML file and control them from a browser: start, stop, restart, tail logs, and open their web interfaces with one click.

## Features

- Start / stop / restart managed processes
- Live log tail (auto-refreshes every 3 seconds)
- One-click browser launch for apps with a web interface
- YAML config — no database, no daemon
- Detects already-running processes on startup (survives restarts)
- Runs on Linux and Windows

## Requirements

- Go 1.24+
- [templ](https://templ.guide) (`go install github.com/a-h/templ/cmd/templ@latest`)
- [Task](https://taskfile.dev) (`go install github.com/go-task/task/v3/cmd/task@latest`) — optional, for the task shortcuts

## Quick start

```sh
git clone <repo>
cd goserviceman
task build          # or: templ generate && go build -o goserviceman .
cp goserviceman.yaml my-apps.yaml   # edit to taste
./goserviceman --config my-apps.yaml
```

Open `http://localhost:9000` in your browser.

## Config file

```yaml
apps:
  - name: "My API"
    executable: "/usr/local/bin/myapi"
    args: ["--port", "8081"]
    port: 8081
    log_file: "/var/log/myapi/app.log"
    start_immediately: true

  - name: "Worker"
    executable: "/usr/local/bin/worker"
    args: []
    port: 0            # 0 = no web interface, Open Browser button is hidden
    log_file: "/var/log/worker/worker.log"
```

| Field        | Description                                              |
|--------------|----------------------------------------------------------|
| `name`       | Display name (must be unique)                            |
| `executable` | Full path to the binary                                  |
| `args`       | Command-line arguments (optional)                        |
| `port`       | Port of the app's web interface; `0` to disable          |
| `log_file`   | Path to the log file to tail                             |
| `start_immediately` | Start the app automatically when goserviceman loads it |

## CLI flags

| Flag       | Default               | Description            |
|------------|-----------------------|------------------------|
| `--config` | `goserviceman.yaml`   | Path to the config file |
| `--addr`   | `:9000`               | Listen address          |

## Tasks

| Command              | Description                                    |
|----------------------|------------------------------------------------|
| `task build`         | Generate templates and build for current OS    |
| `task build-windows` | Cross-compile for Windows (no console window)  |
| `task generate`      | Compile `.templ` files only                    |
| `task run`           | Build and run                                  |

## How it works

goserviceman spawns and tracks processes using the standard `os/exec` package. On startup it scans running processes via [go-ps](https://github.com/mitchellh/go-ps) and reattaches to any that match a configured executable — so restarting goserviceman does not affect managed apps.

The web UI is built with [templ](https://templ.guide), [HTMX](https://htmx.org), and Tailwind CSS. Every action (start, stop, restart, log refresh) is a small HTTP request that swaps only the affected fragment, with no full page reload.

Process termination uses `SIGTERM` on Linux/macOS and `taskkill /F` on Windows.
