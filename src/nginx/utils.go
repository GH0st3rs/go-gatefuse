package nginx

import (
	"os"
)

func ensureDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func writeConf(filePath string, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}
