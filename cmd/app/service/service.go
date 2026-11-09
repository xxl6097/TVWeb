package service

import (
	"fmt"
	"os"

	"github.com/kardianos/service"
	tvbox "github.com/xxl6097/TVBoxWeb"
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

// OnRun starts the TVWeb service. It reads the port from the signed config.
func (t *TVWebService) OnRun(svc igs.Service) error {
	t.gs = svc
	cfg, err := loadConfig()
	if err != nil {
		z.L().Warn("读取配置失败，使用默认端口", zap.Int("port", defaultPort), zap.Error(err))
		cfg = &Config{Port: defaultPort}
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	z.L().Info("Starting TVWeb service", zap.String("addr", addr))

	if err = runServer(addr, svc); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}

func runServer(addr string, gs igs.Service) error {
	upgrader := internal.NewUpgrader(gs)

	opts := []tvbox.Option{
		tvbox.WithDataDir("data"),
		tvbox.WithUpgrader(upgrader),
	}

	if err := tvbox.ListenAndServe(addr, opts...); err != nil {
		z.L().Error("TVBox server error", zap.Error(err))
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
// The returned []byte is embedded into the signed binary; []string are run args.
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
