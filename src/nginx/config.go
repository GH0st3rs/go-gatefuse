package nginx

import (
	"bytes"
	"text/template"
)

type StreamService struct {
	BackendName     string
	ServiceAddress  string
	ServicePort     int
	ExternalPort    int
	ExternalAddress string
}

func resolveTemplate(tmpl string, service interface{}) (string, error) {
	// Parse the template
	t, err := template.New("nginxConfig").Parse(tmpl)
	if err != nil {
		return "", err
	}
	// Execute the template with the struct
	var out bytes.Buffer
	err = t.Execute(&out, service)
	if err != nil {
		return "", err
	}
	// Print the resulting config
	return out.String(), nil
}
