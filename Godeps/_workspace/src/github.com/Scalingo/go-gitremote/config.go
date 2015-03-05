package gitremote

import (
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	Path string
}

func New(path string) *Config {
	configFilePath := filepath.Join(path, ".git", "config")
	config := &Config{Path: configFilePath}
	return config
}

func (config *Config) List() (Remotes, error) {
	content, err := config.Content()
	if err != nil {
		return nil, err
	}
	return ListFromConfigContent(content)
}

func (config *Config) Content() (string, error) {
	content, err := ioutil.ReadFile(config.Path)
	if err != nil {
		return "", &ConfigFileError{err}
	}
	return string(content), nil
}

func (config *Config) Write(content string) error {
	return ioutil.WriteFile(config.Path, []byte(content), 0644)
}
