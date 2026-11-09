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

// Upgrader 实现 tvbox.Upgrader 接口，对接 go-service 的 github 升级能力。
type Upgrader struct {
	gs igs.Service
}

// Check 检查是否有新版本。返回 nil 表示已是最新，返回 Release 表示有新版本可用。
func (u *Upgrader) Check(ctx context.Context) (*tvbox.Release, error) {
	github.Api().SetName(pkg.GithubUser, pkg.GithubRepo)
	data, err := github.Api().CheckUpgrade(pkg.BinName)
	if err != nil {
		// "已是最新版本" 等正常情况也通过 err 返回，统一视为无更新
		z.L().Sugar().Warnln(err)
		return nil, nil
	}
	publishedAt := data["publishedAt"].(time.Time)
	return &tvbox.Release{
		Version:     data["version"].(string),
		Notes:       data["releaseNotes"].(string),
		PublishedAt: publishedAt.Format(time.DateTime),
		PageURL:     data["htmlUrl"].(string),
		FullURL:     data["fullUrl"].(string),
		PatchURL:    data["patchUrl"].(string),
	}, nil
}

// Upgrade 下载并替换二进制，通过 onProgress 回调上报进度。
func (u *Upgrader) Upgrade(ctx context.Context, rel *tvbox.Release, onProgress func(tvbox.UpgradeProgress)) error {
	urls := github.Api().GetProxyUrls(rel.FullURL)

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

	fileUri := utils.DownloadFileWithCancelByUrls(urls, func(progress utils.Progress) {
		report(tvbox.UpgradeProgress{
			Stage:   tvbox.StageDownloading,
			Loaded:  progress.Loaded,
			Total:   progress.Total,
			Message: rel.Version,
		})
	})

	if fileUri == "" {
		z.L().Debug("升级文件下载失败")
		return fmt.Errorf("下载失败，所有镜像地址均不可用")
	}

	z.L().Debug("升级文件下载完成", zap.String("path", fileUri))
	report(tvbox.UpgradeProgress{
		Stage:   tvbox.StageDone,
		Message: "已下载并校验通过，正在替换",
	})
	return u.gs.Upgrade(ctx, fileUri)
}

// NewUpgrader 创建 Upgrader，持有 igs.Service 用于最终的文件替换和重启。
func NewUpgrader(gs igs.Service) *Upgrader {
	return &Upgrader{gs: gs}
}
