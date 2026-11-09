package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/xxl6097/TVWeb/pkg"
	"github.com/xxl6097/glog/pkg/z"
	"go.uber.org/zap"
)

// UpgradeHandler 提供升级相关的 HTTP API
type UpgradeHandler struct {
	upgrader *Upgrader
}

// NewUpgradeHandler 创建升级处理器
func NewUpgradeHandler(upgrader *Upgrader) *UpgradeHandler {
	return &UpgradeHandler{upgrader: upgrader}
}

// CheckUpdateResponse 检查更新的响应
type CheckUpdateResponse struct {
	HasUpdate   bool         `json:"hasUpdate"`
	CurrentVer  string       `json:"currentVersion"`
	LatestVer   string       `json:"latestVersion,omitempty"`
	ReleaseNote *pkg.Release `json:"releaseNote,omitempty"`
	Error       string       `json:"error,omitempty"`
}

// VersionResponse 版本信息响应
type VersionResponse struct {
	AppName      string `json:"appName"`
	AppVersion   string `json:"appVersion"`
	BuildTime    string `json:"buildTime"`
	GitCommit    string `json:"gitCommit"`
	GitBranch    string `json:"gitBranch"`
	OsType       string `json:"osType"`
	Arch         string `json:"arch"`
	GithubUser   string `json:"githubUser"`
	GithubRepo   string `json:"githubRepo"`
}

// HandleCheckUpdate 处理检查更新请求
// GET /api/upgrade/check
func (h *UpgradeHandler) HandleCheckUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	resp := CheckUpdateResponse{
		CurrentVer: pkg.AppVersion,
	}

	rel, err := h.upgrader.Check(ctx)
	if err != nil {
		z.L().Warn("检查更新失败", zap.Error(err))
		resp.Error = err.Error()
	} else if rel != nil && rel.Version != "" && rel.Version != pkg.AppVersion {
		resp.HasUpdate = true
		resp.LatestVer = rel.Version
		resp.ReleaseNote = rel
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleVersion 处理版本信息请求
// GET /api/version
func (h *UpgradeHandler) HandleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := VersionResponse{
		AppName:    pkg.AppName,
		AppVersion: pkg.AppVersion,
		BuildTime:  pkg.BuildTime,
		GitCommit:  pkg.GitCommit,
		GitBranch:  pkg.GitBranch,
		OsType:     pkg.OsType,
		Arch:       pkg.Arch,
		GithubUser: pkg.GithubUser,
		GithubRepo: pkg.GithubRepo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RegisterRoutes 注册路由到 http.ServeMux
func (h *UpgradeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/version", h.HandleVersion)
	mux.HandleFunc("/api/upgrade/check", h.HandleCheckUpdate)
}
