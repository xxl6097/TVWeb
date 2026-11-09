package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/xxl6097/TVWeb/pkg"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/github"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/go-service/pkg/utils"
	"go.uber.org/zap"
)

// Upgrader 实现 pkg.Upgrader。
type Upgrader struct {
	gs igs.Service
}

func (u *Upgrader) Check(ctx context.Context) (*pkg.Release, error) {
	var releaseNote *pkg.Release
	github.Api().SetName(pkg.GithubUser, pkg.GithubRepo)
	data, err := github.Api().CheckUpgrade(pkg.BinName)
	if err != nil {
		// 已是最新版时 err 为 "已是最新版本～"
		z.L().Sugar().Warnln(err)
		return releaseNote, nil
	}
	publishedAt := data["publishedAt"].(time.Time)
	releaseNote = &pkg.Release{
		Version:     data["version"].(string),
		Notes:       data["releaseNotes"].(string),
		PublishedAt: publishedAt.Format(time.DateTime),
		PageURL:     data["htmlUrl"].(string),
		FullURL:     data["fullUrl"].(string),
		PatchURL:    data["patchUrl"].(string),
	}
	return releaseNote, err
}

func (u *Upgrader) Upgrade(ctx context.Context, rel *pkg.Release, onProgress func(pkg.UpgradeProgress)) error {
	urls := github.Api().GetProxyUrls(rel.FullURL)

	report := func(p pkg.UpgradeProgress) {
		if onProgress != nil {
			onProgress(p)
		}
	}

	report(pkg.UpgradeProgress{
		Stage: pkg.StageDownloading, Total: rel.Size, Message: rel.Version,
	})

	fileUri := utils.DownloadFileWithCancelByUrls(urls, func(progress utils.Progress) {
		report(pkg.UpgradeProgress{
			Stage: pkg.StageDownloading, Loaded: progress.Loaded, Total: progress.Total, Message: rel.Version,
		})
	})
	// 并发竞速，返回 "" 表示失败
	if fileUri != "" {
		z.L().Debug("升级文件", zap.String("fileUri", fileUri))
		report(pkg.UpgradeProgress{
			Stage:   pkg.StageDone,
			Message: "已下载并校验通过，准备替换",
		})
		err := u.gs.Upgrade(ctx, fileUri)
		return err
	}

	z.L().Debug("文件地址空fileUri", zap.String("fileUri", fileUri))
	return fmt.Errorf("文件地址空fileUri=%s", fileUri)
}

func NewUpgrader(gs igs.Service) *Upgrader {
	return &Upgrader{gs: gs}
}
