package service

import (
	"context"
	"fmt"
	"os"

	"github.com/kardianos/service"
	"github.com/xxl6097/TVWeb/internal"
	"github.com/xxl6097/TVWeb/pkg"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"go.uber.org/zap"
)

// TVWebService represents the TVWeb service.
type TVWebService struct {
	igs.BaseService
	gs igs.Service
}

// OnConfig returns the service configuration.
func (t *TVWebService) OnConfig() *service.Config {
	return &service.Config{
		Name:        "TVWebService",
		DisplayName: fmt.Sprintf("TVWeb Service %s", pkg.AppVersion),
		Description: "TVWeb - Video Box Web Service",
	}
}

// OnVersion returns the service version.
func (t *TVWebService) OnVersion() string {
	pkg.Version()
	return pkg.AppVersion
}

// OnRun starts the TVWeb service.
func (t *TVWebService) OnRun(svc igs.Service) error {
	t.gs = svc
	cfg, err := loadConfig()
	if err != nil {
		z.L().Warn("读取配置失败，使用默认端口", zap.Int("port", defaultPort), zap.Error(err))
		cfg = &Config{Port: defaultPort}
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	z.L().Info("Starting TVWeb service", zap.String("addr", addr))

	if err = internal.Run(context.Background(), addr, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}

// OnStop stops the TVWeb service.
func (t *TVWebService) OnStop() {
	z.L().Info("Stopping TVWeb service...")
}

// OnShutdown handles service shutdown.
func (t *TVWebService) OnShutdown() {
	z.L().Info("Shutting down TVWeb service...")
}

// GetAny is called during installation to collect configuration from the user.
func (t *TVWebService) GetAny(binDir string) ([]byte, []string) {
	return inputConfig(), nil
}

// OnFinish handles service completion.
func (t *TVWebService) OnFinish() {
	z.L().Info("TVWeb service finished")
}

// UnInstall uninstalls the service.
func (t *TVWebService) UnInstall() {
	z.L().Info("Uninstalling TVWeb service...")
}
