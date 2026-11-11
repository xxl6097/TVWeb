package internal

import (
	"context"
	"fmt"
	"time"

	tvbox "github.com/xxl6097/TVBoxWeb"
	"github.com/xxl6097/TVWeb/pkg"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/github"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/go-service/pkg/utils"
	"go.uber.org/zap"
)

// Upgrader 实现 tvbox.Upgrader 接口。
// 使用 go-service 的 GitHub API 检查更新、代理感知并发下载，
// 最终通过 igs.Service.Upgrade() 完成二进制替换和服务重启。
type Upgrader struct {
	gs igs.Service
}

// Check 检查是否有新版本。nil 表示已是最新，non-nil 表示有更新。
func (u *Upgrader) Check(ctx context.Context) (*tvbox.Release, error) {
	github.Api().SetName(pkg.GithubUser, pkg.GithubRepo)
	data, err := github.Api().CheckUpgrade(pkg.BinName)
	if err != nil {
		// "已是最新版本～" 等正常情况通过 err 返回，视为无更新
		z.L().Sugar().Warnln(err)
		return nil, nil
	}
	publishedAt := data["publishedAt"].(time.Time)
	return &tvbox.Release{
		Version:     data["version"].(string),
		Notes:       data["releaseNotes"].(string),
		PublishedAt: publishedAt.Format(time.DateTime),
		PageURL:     data["htmlUrl"].(string),
		DownloadURL: data["fullUrl"].(string),
	}, nil
}

// Upgrade 下载新版本并通过 go-service 替换二进制、重启服务。
func (u *Upgrader) Upgrade(ctx context.Context, rel *tvbox.Release, onProgress func(tvbox.UpgradeProgress)) error {
	urls := github.Api().GetProxyUrls(rel.DownloadURL)

	report := func(p tvbox.UpgradeProgress) {
		if onProgress != nil {
			onProgress(p)
		}
	}

	report(tvbox.UpgradeProgress{
		Stage:   tvbox.StageDownloading,
		Total:   rel.Size,
		Message: rel.Version,
	})

	// 并发竞速所有镜像地址，返回 "" 表示全部失败
	fileUri := utils.DownloadFileWithCancelByUrls(urls, func(progress utils.Progress) {
		report(tvbox.UpgradeProgress{
			Stage:   tvbox.StageDownloading,
			Loaded:  progress.Loaded,
			Total:   progress.Total,
			Message: rel.Version,
		})
	})

	if fileUri == "" {
		z.L().Debug("升级文件下载失败，所有镜像地址均不可用")
		return fmt.Errorf("下载失败，所有镜像地址均不可用")
	}

	z.L().Debug("升级文件下载完成", zap.String("path", fileUri))
	report(tvbox.UpgradeProgress{
		Stage:   tvbox.StageDone,
		Message: "已下载并校验通过，正在替换",
	})
	// 委托 go-service 完成二进制替换和服务重启
	return u.gs.Upgrade(ctx, fileUri)
}

// NewUpgrader 创建 Upgrader，持有 igs.Service 用于最终的文件替换和重启。
func NewUpgrader(gs igs.Service) *Upgrader {
	return &Upgrader{gs: gs}
}
