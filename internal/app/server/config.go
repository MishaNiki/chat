package server

import (
	"encoding/json"
	"io"
	"os"

	"github.com/MishaNiki/chat/internal/app/templates"
)

// Config ...
type Config struct {
	BindPort  string `json:"bind_port"`
	LogLevel  string `json:"log_level"`
	Templates *templates.Config
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindPort:  ":8085",
		LogLevel:  "debug",
		Templates: templates.NewConfig(),
	}
}

// DecodeJSONConf ...
func (config *Config) DecodeJSONConf(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, 1024)

	var lenBuf int
	for {
		len, e := file.Read(data)
		lenBuf += len
		if e == io.EOF {
			break
		}
	}

	err = json.Unmarshal(data[:lenBuf], config)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data[:lenBuf], config.Templates)
	if err != nil {
		return err
	}

	return nil
}
