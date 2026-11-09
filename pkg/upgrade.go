package pkg

import "context"

// Release 表示一个发布版本的信息
type Release struct {
	Version     string `json:"version"`
	Notes       string `json:"notes"`
	PublishedAt string `json:"publishedAt"`
	PageURL     string `json:"pageUrl"`
	FullURL     string `json:"fullUrl"`
	PatchURL    string `json:"patchUrl"`
	Size        int64  `json:"size"`
}

// UpgradeStage 升级阶段
type UpgradeStage int

const (
	StageDownloading UpgradeStage = iota
	StageVerifying
	StageReplacing
	StageDone
)

// UpgradeProgress 升级进度信息
type UpgradeProgress struct {
	Stage   UpgradeStage
	Loaded  int64
	Total   int64
	Message string
}

// Upgrader 升级器接口
type Upgrader interface {
	// Check 检查是否有新版本
	Check(ctx context.Context) (*Release, error)
	// Upgrade 执行升级
	Upgrade(ctx context.Context, rel *Release, onProgress func(UpgradeProgress)) error
}
