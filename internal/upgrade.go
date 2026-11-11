package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tvbox "github.com/xxl6097/TVBoxWeb"
	"github.com/xxl6097/TVWeb/pkg"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/github"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"go.uber.org/zap"
)

// Upgrader 实现 tvbox.Upgrader 接口。
// Check  使用 go-service github.Api 检查新版本（代理感知）。
// Upgrade 使用带 context/超时的 HTTP 下载，委托 igs.Service.Upgrade() 替换和重启。
type Upgrader struct {
	gs igs.Service
}

// Check 检查是否有新版本。nil 表示已是最新，non-nil 表示有更新可用。
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

// Upgrade 下载新版本并委托 go-service 替换二进制、重启服务。
// 逐一尝试代理地址，每个地址有独立 5 分钟超时，全部失败才报错。
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

	// 逐一尝试，每个地址设独立超时；任一成功即停止
	var filePath string
	var lastErr error
	for _, rawURL := range urls {
		if rawURL == "" {
			continue
		}
		dlCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		path, err := downloadWithProgress(dlCtx, rawURL, rel.Size, func(loaded, total int64) {
			report(tvbox.UpgradeProgress{
				Stage:   tvbox.StageDownloading,
				Loaded:  loaded,
				Total:   total,
				Message: rel.Version,
			})
		})
		cancel()
		if err == nil {
			filePath = path
			break
		}
		lastErr = err
		z.L().Warn("下载失败，尝试下一个地址", zap.String("url", rawURL), zap.Error(err))
	}

	if filePath == "" {
		if lastErr != nil {
			return fmt.Errorf("所有下载地址均失败：%w", lastErr)
		}
		return fmt.Errorf("无可用下载地址")
	}
	defer os.Remove(filePath)

	z.L().Debug("升级文件下载完成", zap.String("path", filePath))
	report(tvbox.UpgradeProgress{
		Stage:   tvbox.StageReplacing,
		Message: "已下载，正在替换",
	})

	// 委托 go-service 完成二进制替换和服务重启
	return u.gs.Upgrade(ctx, filePath)
}

// downloadWithProgress 通过 HTTP GET 下载 url 到临时文件，每 200ms 上报一次进度。
// ctx 可用于取消或超时控制（调用方应设置 WithTimeout）。
func downloadWithProgress(ctx context.Context, rawURL string, hint int64, onProgress func(loaded, total int64)) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "TVWeb-upgrade")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	total := resp.ContentLength
	if total <= 0 {
		total = hint // 服务端未返回 Content-Length 时用 release 里的 size 兜底
	}

	tmp, err := os.CreateTemp("", "tvweb-upgrade-*")
	if err != nil {
		return "", err
	}
	tmpName := tmp.Name()

	var written int64
	buf := make([]byte, 32*1024)
	lastReport := time.Now()

	for {
		select {
		case <-ctx.Done():
			tmp.Close()
			os.Remove(tmpName)
			return "", ctx.Err()
		default:
		}

		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			nw, ew := tmp.Write(buf[:nr])
			written += int64(nw)
			if ew != nil {
				tmp.Close()
				os.Remove(tmpName)
				return "", ew
			}
			if time.Since(lastReport) >= 200*time.Millisecond {
				onProgress(written, total)
				lastReport = time.Now()
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			tmp.Close()
			os.Remove(tmpName)
			return "", er
		}
	}
	tmp.Close()

	// 最后一次完整进度上报
	onProgress(written, total)

	// 确保可执行
	_ = os.Chmod(tmpName, 0755)
	return tmpName, nil
}

// NewUpgrader 创建 Upgrader，持有 igs.Service 用于最终的文件替换和重启。
func NewUpgrader(gs igs.Service) *Upgrader {
	return &Upgrader{gs: gs}
}
