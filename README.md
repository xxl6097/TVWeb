# TVWeb

TVWeb 是一个基于 Go 语言的视频盒子 Web 服务框架，集成了 TVBoxWeb 的功能，提供了完整的服务管理和日志记录能力。

## 项目特性

- 🚀 基于 `go-service` 框架的服务管理
- 📦 集成 TVBoxWeb 的核心功能
- 📝 完整的日志记录系统
- 🛠️ 模块化的项目结构
- ⚙️ 自动升级和重启机制

## 快速开始

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

- `github.com/xxl6097/TVBoxWeb` - TVBox Web 框架
- `github.com/xxl6097/glog` - 日志库
- `github.com/xxl6097/go-service` - 服务管理框架
- `go.uber.org/zap` - 高性能日志记录

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
