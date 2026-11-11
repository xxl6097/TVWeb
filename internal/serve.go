package internal

import (
	"context"
	"fmt"
	"io"

	tvbox "github.com/xxl6097/TVBoxWeb"
	"github.com/xxl6097/TVBoxWeb/githubupgrade"
	"github.com/xxl6097/TVWeb/pkg"
	"github.com/xxl6097/glog/pkg/z"
	"go.uber.org/zap"
)

// Run 启动 TVWeb 服务，参考 go-fs 的 RunCfg 写法：
//   - ctx 取消时优雅退出
//   - 服务自行出错时返回错误
func Run(ctx context.Context, addr string, stdout, stderr io.Writer) error {
	opts := []tvbox.Option{
		tvbox.WithDataDir("data"),
		// 把 TVWeb 版本号写入 TVBoxWeb 内部，供升级检测和界面展示使用
		tvbox.WithVersion(pkg.AppVersion),
	}

	// 注入 GitHub 升级器（v1.6.0 内置：SHA256 校验 + 断点续传 + 双通道下载）
	up, err := githubupgrade.New(githubupgrade.Options{
		Owner: pkg.GithubUser,
		Repo:  pkg.GithubRepo,
	})
	if err != nil {
		fmt.Fprintf(stderr, "TVWeb: 在线升级未启用: %v\n", err)
		z.L().Warn("在线升级未启用", zap.Error(err))
	} else {
		opts = append(opts, tvbox.WithUpgrader(up))
		z.L().Info("在线升级已启用", zap.String("repo", pkg.GithubUser+"/"+pkg.GithubRepo))
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

	fmt.Fprintf(stdout, "TVWeb 服务已启动：http://localhost%s\n", addr)

	select {
	case <-ctx.Done():
		fmt.Fprintln(stderr, "\nTVWeb: 正在关闭…")
		return nil
	case err := <-done:
		return err
	}
}
