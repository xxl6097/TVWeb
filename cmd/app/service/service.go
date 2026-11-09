package service

import (
	"flag"
	"fmt"
	"os"

	"github.com/kardianos/service"
	tvbox "github.com/xxl6097/TVBoxWeb"
	"github.com/xxl6097/TVWeb/pkg"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"go.uber.org/zap"
)

// TVWebService represents the TVWeb service.
type TVWebService struct{}

// OnConfig returns the service configuration.
func (t *TVWebService) OnConfig() *service.Config {
	return &service.Config{
		Name:        "TVWebService",
		DisplayName: "TVWeb Service",
		Description: "TVWeb - Video Box Web Service",
	}
}

// OnVersion returns the service version.
func (t *TVWebService) OnVersion() string {
	return pkg.AppVersion
}

// OnRun starts the TVWeb service. It reads the port from the signed config.
func (t *TVWebService) OnRun(svc igs.Service) error {
	cfg, err := loadConfig()
	if err != nil {
		z.L().Warn("读取配置失败，使用默认端口", zap.Int("port", defaultPort), zap.Error(err))
		cfg = &Config{Port: defaultPort}
	}

	z.L().Info("Starting TVWeb service", zap.Int("port", cfg.Port))

	// TODO: pass cfg.Port to TVBox when integrated
	// e.g.: tvbox.Start(cfg.Port)

	addr := fmt.Sprintf(":%d", cfg.Port)
	dataDir := "data"
	flag.Parse()

	if err = tvbox.ListenAndServe(addr, tvbox.WithDataDir(dataDir)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	z.L().Info("TVWeb service started successfully", zap.String("addr", addr))
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
