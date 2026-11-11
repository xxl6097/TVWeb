package service

import (
	"fmt"
	"os"

	"github.com/kardianos/service"
	tvbox "github.com/xxl6097/TVBoxWeb"
	"github.com/xxl6097/TVBoxWeb/githubupgrade"
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

	if err = runServer(addr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}

func runServer(addr string) error {
	opts := []tvbox.Option{
		tvbox.WithDataDir("data"),
		// 把 TVWeb 的版本号注入到 TVBoxWeb 的 api.AppVersion，
		// 供升级检测和界面展示使用
		tvbox.WithVersion(pkg.AppVersion),
	}

	// 使用 TVBoxWeb v1.6.0 内置的 GitHub 升级器；
	// 检查 github.com/xxl6097/TVWeb 的 Release，
	// 查找名为 yezhen-<os>-<arch> 的资产文件（含 checksums.txt 校验）。
	up, err := githubupgrade.New(githubupgrade.Options{
		Owner: pkg.GithubUser,
		Repo:  pkg.GithubRepo,
	})
	if err != nil {
		z.L().Warn("在线升级未启用", zap.Error(err))
	} else {
		opts = append(opts, tvbox.WithUpgrader(up))
		z.L().Info("在线升级已启用",
			zap.String("repo", pkg.GithubUser+"/"+pkg.GithubRepo))
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
