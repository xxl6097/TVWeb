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
	User string // 登录用户名（空表示不启用登录保护）
	Pass string // 登录密码
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
	user := utils.InputStringEmpty("请输入登录用户名（留空不启用登录保护）：", "")
	var pass string
	if user != "" {
		pass = utils.InputString("请输入登录密码：")
	}
	cfg := Config{Port: port, User: user, Pass: pass}
	bb, err := ukey.StructToGob(cfg)
	if err != nil {
		return nil
	}
	return bb
}
