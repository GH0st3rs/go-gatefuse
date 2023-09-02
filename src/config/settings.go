package config

import (
	"flag"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
)

type AppSettings struct {
	MainDomain    string `json:"main_domain" form:"main_domain"`
	NginxConfPath string `json:"nginx_conf_path" form:"nginx_conf_path"`
	Username      string `json:"username" form:"username"`
	Password      string `json:"password" form:"password"`
	RequestType   string `json:"request_type" form:"request_type"`
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
	AppPort  = *flag.Int("app_port", 3000, "Listening port")
	AppHost  = *flag.String("app_host", "0.0.0.0", "Listening address")
	AppDebug = *flag.Bool("debug", true, "Enable debug print")
	Settings AppSettings
	// Fiber Variables
	SqliteStorage  = sqlite3.New()
	SessionStorage = session.New(session.Config{
		Storage: SqliteStorage,
	})
)
