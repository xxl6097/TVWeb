package service

import (
	"fmt"

	"github.com/xxl6097/go-service/pkg/ukey"
	"github.com/xxl6097/go-service/pkg/utils"
)

const defaultPort = 8888

// Config holds the service configuration persisted in the signed binary.
type Config struct {
	Port int
}

// loadConfig reads the Config embedded in the signed binary via ukey.
func loadConfig() (*Config, error) {
	buf, err := ukey.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	var cfg Config
	if err := ukey.GobToStruct(buf, &cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	return &cfg, nil
}

// inputConfig prompts the user for configuration values during installation.
func inputConfig() []byte {
	port := utils.InputIntDefault(fmt.Sprintf("请输入端口号（默认 %d）：", defaultPort), defaultPort)
	cfg := Config{Port: port}
	bb, err := ukey.StructToGob(cfg)
	if err != nil {
		return nil
	}
	return bb
}
