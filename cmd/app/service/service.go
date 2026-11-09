package service

import (
	"context"
	"flag"
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

	z.L().Info("Starting TVWeb service", zap.Int("port", cfg.Port))

	addr := fmt.Sprintf(":%d", cfg.Port)
	dataDir := "data"
	flag.Parse()

	if err = runServer(addr, dataDir, svc); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	z.L().Info("TVWeb service started successfully", zap.String("addr", addr))
	return nil
}

func runServer(addr, dataDir string, gs igs.Service) error {
	// 初始化升级器（后台定期检查更新）
	upgrader := internal.NewUpgrader(gs)

	// 启动时检查一次更新
	go func() {
		ctx := context.Background()
		if rel, err := upgrader.Check(ctx); err == nil && rel != nil {
			z.L().Info("发现新版本", zap.String("version", rel.Version))
		}
	}()

	if err := tvbox.ListenAndServe(addr, tvbox.WithDataDir(dataDir)); err != nil {
		z.L().Error("Failed to start TVBox server", zap.Error(err))
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
