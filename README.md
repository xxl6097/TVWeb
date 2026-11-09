# TVWeb

[![Build and Test](https://github.com/xxl6097/TVWeb/actions/workflows/build.yml/badge.svg)](https://github.com/xxl6097/TVWeb/actions/workflows/build.yml)
[![Release](https://github.com/xxl6097/TVWeb/actions/workflows/release.yml/badge.svg)](https://github.com/xxl6097/TVWeb/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/xxl6097/TVWeb)](https://goreportcard.com/report/github.com/xxl6097/TVWeb)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/xxl6097/TVWeb)](https://github.com/xxl6097/TVWeb)
[![Release](https://img.shields.io/github/v/release/xxl6097/TVWeb)](https://github.com/xxl6097/TVWeb/releases)

TVWeb 是一个基于 Go 语言的视频盒子 Web 服务框架，集成了 TVBoxWeb 的功能，提供了完整的服务管理和日志记录能力。

## 项目特性

- 🚀 基于 `go-service` 框架的服务管理
- 📦 集成 TVBoxWeb 的核心功能
- 📝 完整的日志记录系统
- 🛠️ 模块化的项目结构
- ⚙️ 自动升级和重启机制

## 快速开始

### 下载预编译二进制

从 [Releases](https://github.com/xxl6097/TVWeb/releases) 页面下载适合你系统的二进制文件：

**Linux:**
```bash
wget https://github.com/xxl6097/TVWeb/releases/latest/download/TVWeb-linux-amd64
chmod +x TVWeb-linux-amd64
./TVWeb-linux-amd64
```

**macOS:**
```bash
# Intel 芯片
wget https://github.com/xxl6097/TVWeb/releases/latest/download/TVWeb-darwin-amd64
chmod +x TVWeb-darwin-amd64
./TVWeb-darwin-amd64

# Apple Silicon (M1/M2)
wget https://github.com/xxl6097/TVWeb/releases/latest/download/TVWeb-darwin-arm64
chmod +x TVWeb-darwin-arm64
./TVWeb-darwin-arm64
```

**Windows:**

从 Releases 下载 `TVWeb-windows-amd64.exe` 并直接运行。

### 从源码构建

### 前置需求

- Go 1.26 或更高版本
- Git

### 安装

```bash
git clone https://github.com/xxl6097/TVWeb.git
cd TVWeb
go mod download
```

### 编译

```bash
go build -o TVWeb cmd/app/main.go
```

### 运行

```bash
./TVWeb
```

## 项目结构

```
TVWeb/
├── .github/
│   └── workflows/          # GitHub Actions 工作流
│       ├── build.yml       # 构建和测试
│       ├── release.yml     # 自动发布
│       ├── pr-check.yml    # PR 检查
│       └── manual-build.yml # 手动构建
├── cmd/
│   └── app/
│       ├── main.go          # 应用入口
│       └── service/
│           └── service.go   # 服务实现
├── internal/                # 内部实现
├── pkg/
│   └── version.go          # 版本信息
├── go.mod                   # 依赖配置
└── README.md               # 项目文档
```

## 依赖

- `github.com/kardianos/service` - 跨平台服务管理
- `github.com/xxl6097/glog` - 日志库
- `github.com/xxl6097/go-service` - 服务管理框架
- `go.uber.org/zap` - 高性能日志记录

## CI/CD

本项目使用 GitHub Actions 进行自动化构建和发布：

- **构建和测试** (`build.yml`): 每次推送到 main/develop 分支时自动运行
- **发布** (`release.yml`): 创建新 tag 时自动构建多平台二进制并发布到 Releases
- **PR 检查** (`pr-check.yml`): 对 Pull Request 进行代码检查、测试和构建
- **手动构建** (`manual-build.yml`): 支持手动触发，可自定义构建平台和架构

### 发布新版本

```bash
# 创建并推送 tag
git tag -a v0.0.2 -m "Release v0.0.2"
git push origin v0.0.2

# GitHub Actions 将自动构建并发布到 Releases
```

## 开发

### 获取依赖

```bash
go mod tidy
```

### 运行测试

```bash
go test ./...
```

## 许可证

MIT

## 作者

xxl6097

## 相关项目

- [go-fs](https://github.com/xxl6097/go-fs) - 参考项目
- [TVBoxWeb](https://github.com/xxl6097/TVBoxWeb) - 依赖项目
