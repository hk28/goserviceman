package views

// PageVM is the view model for the full dashboard page.
type PageVM struct {
	Apps []AppVM
}

// AppVM is the view model for a single managed application row.
type AppVM struct {
	Name       string
	Executable string
	LogFile    string
	Port       int
	Running    bool
	StatusText string
}

// LogVM is the view model for a log-tail fragment.
type LogVM struct {
	AppName string
	Lines   []string
	Error   string // non-empty when the log file could not be read
}
