package internal

import (
	"context"
	"fmt"
	"io"

	tvbox "github.com/xxl6097/TVBoxWeb"
	"github.com/xxl6097/TVWeb/pkg"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"go.uber.org/zap"
)

// Run 启动 TVWeb 服务，参考 go-fs 的 RunCfg 写法：
//   - ctx 取消时优雅退出
//   - 服务自行出错时返回错误
func Run(ctx context.Context, addr, user, pass string, gs igs.Service, stdout, stderr io.Writer) error {
	opts := []tvbox.Option{
		tvbox.WithDataDir("data"),
		// 把 TVWeb 版本号写入 TVBoxWeb 内部，供升级检测和界面展示使用
		tvbox.WithVersion(pkg.AppVersion),
		// 注入升级器：go-service 代理感知下载 + 委托 igs.Service 替换和重启
		tvbox.WithUpgrader(NewUpgrader(gs)),
	}

	// 有用户名密码时启用登录保护（v1.7.0 新增）
	if user != "" && pass != "" {
		opts = append(opts, tvbox.WithLogin(user, pass))
		z.L().Info("登录保护已启用", zap.String("user", user))
	}

	srv, err := tvbox.New(opts...)
	if err != nil {
		return err
	}
	defer srv.Close()

	// tvbox.Server 没有 Start/Done/Shutdown，用 goroutine + channel 模拟
	done := make(chan error, 1)
	go func() {
		done <- srv.ListenAndServe(addr)
	}()

	z.L().Info("在线升级已启用", zap.String("repo", pkg.GithubUser+"/"+pkg.GithubRepo))
	fmt.Fprintf(stdout, "TVWeb 服务已启动：http://localhost%s\n", addr)

	select {
	case <-ctx.Done():
		fmt.Fprintln(stderr, "\nTVWeb: 正在关闭…")
		return nil
	case err := <-done:
		return err
	}
}
