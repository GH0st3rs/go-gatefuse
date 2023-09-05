package nginx

import (
	"errors"
	"fmt"
	"go-gatefuse/src/config"
	"os"
	"os/exec"
)

func SaveNginxConfig(record config.GateRecord) error {
	service := StreamService{
		BackendName:     record.UUID,
		ServiceAddress:  record.SrcAddr,
		ServicePort:     record.SrcPort,
		ExternalPort:    record.DstPort,
		ExternalAddress: record.DstAddr,
	}

	switch record.Protocol {
	case "udp":
		return CreateUdpConfiguration(service, record.Active)
	case "tcp":
		return CreateTcpConfiguration(service, record.Active)
	}
	return nil
}

func ToggleConfig(record config.GateRecord) error {
	if record.Active {
		if err := SaveNginxConfig(record); err != nil {
			return err
		}
	} else {
		conf_file := fmt.Sprintf("/etc/nginx/conf.d/%s.conf", record.UUID)
		if err := os.Remove(conf_file); err != nil {
			return err
		}
	}
	return ReloadNginx()
}

func ReloadNginx() error {
	_, err := exec.Command("nginx", "-t").Output()
	var e *exec.ExitError
	if err != nil && errors.As(err, &e) {
		return fmt.Errorf("%s", e.Stderr)
	}
	_, err = exec.Command("nginx", "-s", "reload").Output()
	if err != nil && errors.As(err, &e) {
		return fmt.Errorf("%s", e.Stderr)
	}
	return nil
}

func DeleteNginxConfig(record config.GateRecord) error {
	service := StreamService{
		BackendName:     record.UUID,
		ServiceAddress:  record.SrcAddr,
		ServicePort:     record.SrcPort,
		ExternalPort:    record.DstPort,
		ExternalAddress: record.DstAddr,
	}
	switch record.Protocol {
	case "udp":
		return DeleteUdpConfiguration(service)
	case "tcp":
		return DeleteTcpConfiguration(service)
	}
	return nil
}
