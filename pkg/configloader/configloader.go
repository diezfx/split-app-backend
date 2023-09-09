package configloader

import (
	"fmt"
	"os"
)

type Source string

const (
	EnvSource  Source = "Env"
	FileSource Source = "File"
)

type Loader struct {
	Source
	configPath  string
	secretsPath string
}

func NewFileLoader(configPath, secretsPath string) *Loader {
	return &Loader{configPath: configPath, secretsPath: secretsPath}
}

func (c *Loader) LoadConfig(namespace string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s.json", c.configPath, namespace)
	content, err := os.ReadFile(fmt.Sprintf("%s/%s.json", c.configPath, namespace))
	if err != nil {
		return nil, fmt.Errorf("read config file %s: %w", path, err)
	}
	return content, nil
}

func (c *Loader) LoadSecret(namespace, key string) (string, error) {
	path := fmt.Sprintf("%s/%s/%s", c.configPath, namespace, key)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read secret file %s: %w", path, err)
	}
	return string(content), nil
}
