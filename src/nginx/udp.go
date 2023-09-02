package nginx

import (
	"fmt"
	"go-gatefuse/src/config"
	"os"
	"path/filepath"
)

const UDP_PROXY_TEMPLATE string = `
    upstream {{.BackendName}} {
        server {{.ServiceAddress}}:{{.ServicePort}};
    }

    server {
        listen {{.ExternalPort}} udp;
        proxy_pass {{.BackendName}};
        proxy_timeout 1s;
        proxy_responses 1;
        error_log /var/log/nginx/{{.BackendName}}_error.log;
    }`

// CreateUdpConfiguration tried to create *.conf files inside nginx
// directory to allow to use stream UDP reverse-proxy mode
func CreateUdpConfiguration(service StreamService, active bool) error {
	templ, err := resolveTemplate(UDP_PROXY_TEMPLATE, service)
	if err != nil {
		return err
	}

	if active {
		// Create a file if record is active
		conf_file := fmt.Sprintf("%s.conf", service.BackendName)
		conf_file = filepath.Join(config.Settings.NginxConfPath, conf_file)
		return writeConf(conf_file, templ)
	}

	return nil
}

// DeleteUdpConfiguration Delete configuration files for UDP mode
func DeleteUdpConfiguration(service StreamService) error {
	conf_file := fmt.Sprintf("%s.conf", service.BackendName)
	conf_file = filepath.Join(config.Settings.NginxConfPath, conf_file)
	return os.Remove(conf_file)
}
