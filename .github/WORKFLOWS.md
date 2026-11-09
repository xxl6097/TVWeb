# GitHub Actions Workflows

本项目使用 GitHub Actions 实现自动化构建、测试和发布。

## 工作流说明

### 1. Build and Test (`build.yml`)

**触发条件：**
- 推送到 `main` 或 `develop` 分支
- 创建针对 `main` 或 `develop` 分支的 Pull Request

**功能：**
- 代码格式检查 (`gofmt`)
- 静态代码分析 (`go vet`)
- 运行单元测试
- 生成代码覆盖率报告
- 构建二进制文件

**查看状态：**
```bash
gh run list --workflow=build.yml
```

### 2. Release (`release.yml`)

**触发条件：**
- 推送以 `v` 开头的 tag（例如：`v0.0.1`，`v1.0.0`）

**功能：**
- 自动构建多平台二进制文件：
  - Linux (AMD64, ARM64)
  - macOS (AMD64, ARM64)
  - Windows (AMD64, ARM64)
- 生成 SHA256 校验和
- 创建 GitHub Release
- 自动上传所有构建产物

**使用方法：**
```bash
# 1. 更新版本号
vim pkg/version.go  # 修改版本号

# 2. 提交变更
git add pkg/version.go
git commit -m "Bump version to v0.0.2"
git push

# 3. 创建并推送 tag
git tag -a v0.0.2 -m "Release v0.0.2"
git push origin v0.0.2

# 4. GitHub Actions 将自动构建并发布
```

**查看发布状态：**
```bash
gh run list --workflow=release.yml
```

### 3. PR Check (`pr-check.yml`)

**触发条件：**
- 创建、同步或重新打开 Pull Request

**功能：**
- 代码质量检查（golangci-lint）
- 运行所有测试
- 多平台构建验证
- 代码覆盖率报告

**检查内容：**
- Lint: 代码风格和最佳实践检查
- Test: 单元测试和竞态检测
- Build: 确保所有平台都能成功构建

### 4. Manual Build (`manual-build.yml`)

**触发条件：**
- 手动触发（通过 GitHub UI 或 API）

**功能：**
- 自定义平台构建（Linux, macOS, Windows）
- 自定义架构（AMD64, ARM64）
- 构建产物保留 30 天

**使用方法：**

1. 通过 GitHub UI：
   - 访问 `Actions` 标签
   - 选择 `Manual Build` workflow
   - 点击 `Run workflow`
   - 选择分支和构建选项
   - 点击 `Run workflow`

2. 通过 GitHub CLI：
```bash
gh workflow run manual-build.yml \
  -f platforms=linux,darwin \
  -f architectures=amd64,arm64
```

**下载构建产物：**
```bash
# 列出最近的 workflow 运行
gh run list --workflow=manual-build.yml

# 下载产物
gh run download <run-id>
```

## 构建平台矩阵

| 操作系统 | 架构 | 文件名 | 支持的 Workflows |
|---------|------|--------|-----------------|
| Linux | AMD64 | TVWeb-linux-amd64 | 所有 |
| Linux | ARM64 | TVWeb-linux-arm64 | Release, Manual |
| macOS | AMD64 | TVWeb-darwin-amd64 | 所有 |
| macOS | ARM64 | TVWeb-darwin-arm64 | 所有 |
| Windows | AMD64 | TVWeb-windows-amd64.exe | 所有 |
| Windows | ARM64 | TVWeb-windows-arm64.exe | Release, Manual |

## 环境变量

所有 workflows 使用的关键环境变量：

- `CGO_ENABLED=0`: 禁用 CGO，确保静态链接
- `GOOS`: 目标操作系统
- `GOARCH`: 目标架构
- `LDFLAGS`: 链接器标志，用于嵌入版本信息

## 版本信息嵌入

构建时会自动嵌入以下信息到二进制文件：

```go
github.com/xxl6097/TVWeb/pkg/version.Version     // 从 tag 提取
github.com/xxl6097/TVWeb/pkg/version.GitCommit   // Git commit hash
github.com/xxl6097/TVWeb/pkg/version.BuildTime   // 构建时间
```

查看二进制文件版本：
```bash
./TVWeb-linux-amd64 --version  # 假设实现了 --version 标志
```

## 故障排查

### Workflow 失败

1. 查看详细日志：
```bash
gh run view <run-id> --log
```

2. 重新运行失败的 job：
```bash
gh run rerun <run-id> --failed
```

### 构建失败

常见问题：
- **依赖问题**: 确保 `go.mod` 和 `go.sum` 已提交
- **测试失败**: 本地运行 `go test ./...` 验证
- **Lint 错误**: 本地运行 `golangci-lint run`

### Release 失败

1. 检查 tag 格式：必须以 `v` 开头
2. 确保有 `contents: write` 权限
3. 检查是否有同名 release 已存在

## 权限配置

Workflows 需要以下权限：

- `contents: write` - Release workflow（创建 release）
- `actions: read` - 所有 workflows（访问 artifacts）

这些权限在 workflow 文件中已配置，无需额外设置。

## 最佳实践

1. **版本管理**
   - 使用语义化版本（Semantic Versioning）
   - Tag 格式：`v{major}.{minor}.{patch}`

2. **发布流程**
   - 在 `develop` 分支开发
   - 合并到 `main` 分支
   - 从 `main` 分支创建 release tag

3. **测试**
   - 推送前确保本地测试通过
   - PR 会自动运行完整测试套件

4. **构建验证**
   - 使用 Manual Build 测试特定平台
   - Release 前先在 PR 中验证构建

## 相关链接

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Go 语言构建最佳实践](https://go.dev/doc/build)
- [golangci-lint 配置](https://golangci-lint.run/usage/configuration/)
