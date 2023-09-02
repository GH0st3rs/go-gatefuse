package nginx

import (
	"fmt"
	"go-gatefuse/src/config"
	"os"
	"path/filepath"
)

const (
	TCP_PROXY_TEMPLATE string = `
	map $ssl_preread_server_name $port{{.ExternalPort}} {
		include /etc/nginx/conf.d/{{.ExternalPort}}/mappings/*.conf;
	}

	include /etc/nginx/conf.d/{{.ExternalPort}}/backends/*.conf;

	server {
		listen      {{.ExternalPort}};
		proxy_pass  $port{{.ExternalPort}};
		ssl_preread on;
	}`

	TCP_BACKEND_TEMPLATE string = `    upstream {{.BackendName}} {
        server {{.ServiceAddress}}:{{.ServicePort}};
    }`

	TCP_MAP_TEMPLATE string = `		{{.ExternalAddress}}      {{.BackendName}};`
)

// CreateTcpConfiguration tried to create *.conf files inside nginx
// directory to allow to use stream TCP reverse-proxy mode
func CreateTcpConfiguration(service StreamService, active bool) error {
	baseDir := fmt.Sprintf("%d", service.ExternalPort)
	baseDir = filepath.Join(config.Settings.NginxConfPath, baseDir)
	mappingDir := filepath.Join(baseDir, "mappings")
	backendDir := filepath.Join(baseDir, "backends")

	if active {
		if err := createMapConf(mappingDir, service); err != nil {
			return err
		}
		if err := createBackendConf(backendDir, service); err != nil {
			return err
		}
		if err := createPortConf(baseDir, service); err != nil {
			return err
		}
	}

	return nil
}

func createPortConf(baseDir string, service StreamService) error {
	templ, err := resolveTemplate(TCP_PROXY_TEMPLATE, service)
	if err != nil {
		return err
	}
	// Create port-specific conf file
	portConf := fmt.Sprintf("%s.conf", baseDir)
	if _, err := os.Stat(portConf); os.IsNotExist(err) {
		if err := writeConf(portConf, templ); err != nil {
			return err
		}
	}
	return nil
}

func createMapConf(mappingDir string, service StreamService) error {
	templ, err := resolveTemplate(TCP_MAP_TEMPLATE, service)
	if err != nil {
		return err
	}
	//Create mappings folder if needed
	if err := ensureDirectory(mappingDir); err != nil {
		return err
	}
	// Create mappings conf file
	mappingConf := fmt.Sprintf("%s/%s.conf", mappingDir, service.BackendName)
	if err := writeConf(mappingConf, templ); err != nil {
		return err
	}
	return nil
}

func createBackendConf(backendDir string, service StreamService) error {
	templ, err := resolveTemplate(TCP_BACKEND_TEMPLATE, service)
	if err != nil {
		return err
	}
	//Create backends folder if needed
	if err := ensureDirectory(backendDir); err != nil {
		return err
	}
	// Create backends conf file
	backendConf := fmt.Sprintf("%s/%s.conf", backendDir, service.BackendName)
	if err := writeConf(backendConf, templ); err != nil {
		return err
	}
	return nil
}
