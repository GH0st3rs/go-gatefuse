package config

import (
	"flag"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
)

type AppSettings struct {
	// Main App settings
	MainDomain        string `json:"main_domain" form:"main_domain"`
	NginxConfPath     string `json:"nginx_conf_path" form:"nginx_conf_path"`
	UnboundConfPath   string `json:"unbound_conf_path" form:"unbound_conf_path"`
	UnboundRemote     bool   `json:"unbound_remote" form:"unbound_remote"`
	UnboundRemoteHost string `json:"unbound_remote_host" form:"unbound_remote_host"`
	// Authentication settings
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type GateRecord struct {
	SrcPort  int    `json:"SrcPort"`
	SrcAddr  string `json:"SrcAddr"`
	DstPort  int    `json:"DstPort"`
	DstAddr  string `json:"DstAddr"`
	Protocol string `json:"Protocol"`
	Comment  string `json:"Comment"`
	Active   bool   `json:"Active"`
	UUID     string `json:"UUID"`
}

var (
	// App Variables
	AppPort  = flag.Int("app_port", 3000, "Listening port")
	AppHost  = flag.String("app_host", "0.0.0.0", "Listening address")
	AppDebug = flag.Bool("debug", true, "Enable debug print")
	UseCache = flag.Bool("cache", false, "Use web cache")
	AppInit  = flag.Bool("init", false, "Initialize (or re-initialize) the database")
	Settings AppSettings
	// Fiber Variables
	SqliteStorage  = sqlite3.New()
	SessionStorage = session.New(session.Config{
		Storage: SqliteStorage,
	})
)
