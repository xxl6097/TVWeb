# Contributing to TVWeb

感谢你对 TVWeb 项目的关注！我们欢迎各种形式的贡献。

## 开发环境设置

### 前置需求

- Go 1.26 或更高版本
- Git
- golangci-lint（可选，用于本地 lint 检查）

### 克隆仓库

```bash
git clone https://github.com/xxl6097/TVWeb.git
cd TVWeb
```

### 安装依赖

```bash
go mod download
```

### 构建

```bash
go build -o TVWeb cmd/app/main.go
```

## 开发流程

### 1. 创建分支

```bash
git checkout -b feature/your-feature-name
# 或
git checkout -b fix/your-bug-fix
```

分支命名规范：
- `feature/` - 新功能
- `fix/` - Bug 修复
- `docs/` - 文档更新
- `refactor/` - 代码重构
- `test/` - 测试相关

### 2. 编写代码

- 遵循 Go 语言编码规范
- 添加必要的注释
- 确保代码通过 lint 检查
- 添加或更新相关测试

### 3. 本地测试

运行所有测试：
```bash
go test ./...
```

运行带覆盖率的测试：
```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 4. 代码格式化

```bash
go fmt ./...
go vet ./...
```

### 5. Lint 检查（可选）

安装 golangci-lint：
```bash
# macOS
brew install golangci-lint

# Linux/Windows
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

运行 lint：
```bash
golangci-lint run
```

### 6. 提交变更

提交信息格式：
```
<type>: <subject>

<body>

<footer>
```

类型（type）：
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响代码运行的变动）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

示例：
```bash
git commit -m "feat: add new API endpoint for user management

- Implement user CRUD operations
- Add validation middleware
- Update API documentation

Closes #123"
```

### 7. 推送分支

```bash
git push origin feature/your-feature-name
```

### 8. 创建 Pull Request

1. 访问 GitHub 仓库页面
2. 点击 "New Pull Request"
3. 选择你的分支
4. 填写 PR 描述：
   - 变更内容概述
   - 相关 Issue 编号
   - 测试说明
   - 截图（如果适用）

## Pull Request 指南

### PR 标题

使用与 commit 相同的格式：
```
feat: add user authentication
fix: resolve memory leak in service loop
docs: update installation instructions
```

### PR 描述模板

```markdown
## 变更类型
- [ ] 新功能
- [ ] Bug 修复
- [ ] 文档更新
- [ ] 代码重构
- [ ] 其他

## 变更说明
简要描述本次变更的内容和目的。

## 相关 Issue
Closes #<issue-number>

## 测试说明
- [ ] 添加了新的测试
- [ ] 所有测试通过
- [ ] 本地测试通过

## 截图（如适用）
附上相关截图

## 检查清单
- [ ] 代码遵循项目编码规范
- [ ] 已运行 `go fmt` 和 `go vet`
- [ ] 已添加或更新相关测试
- [ ] 已更新相关文档
- [ ] 所有 CI 检查通过
```

## 代码规范

### Go 代码风格

遵循 [Effective Go](https://go.dev/doc/effective_go) 和 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)。

关键要点：
- 使用 `gofmt` 格式化代码
- 导出的函数和类型必须有注释
- 错误处理要明确
- 使用有意义的变量名
- 保持函数简短和专注

### 注释规范

```go
// PackageName does something useful.
package main

// FunctionName performs a specific task.
// It returns an error if something goes wrong.
func FunctionName() error {
    // 实现细节
}
```

### 错误处理

```go
// Good
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad
if err != nil {
    return err  // 没有上下文信息
}
```

## 测试指南

### 编写测试

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "test",
            want:    "expected",
            wantErr: false,
        },
        // 更多测试案例...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("FunctionName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 测试覆盖率

目标：保持测试覆盖率在 70% 以上

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## 文档

### 更新文档

如果你的变更涉及：
- API 变化：更新相关文档
- 新功能：在 README.md 中添加使用说明
- 配置变更：更新配置文档

### 文档风格

- 使用清晰、简洁的语言
- 提供代码示例
- 使用 Markdown 格式
- 添加目录（对长文档）

## 发布流程

核心维护者负责发布新版本：

1. 更新版本号
2. 更新 CHANGELOG
3. 创建 release tag
4. GitHub Actions 自动构建和发布

## 行为准则

### 基本原则

- 尊重所有贡献者
- 欢迎建设性的批评
- 专注于对项目最有利的事情
- 表现出对社区成员的同理心

### 不可接受的行为

- 使用性化语言或意象
- 人身攻击
- 骚扰行为
- 公开或私下的骚扰
- 未经许可发布他人私人信息

## 获得帮助

如果你需要帮助：

1. 查看现有文档和 Issue
2. 在 Issue 中提问
3. 联系维护者

## 许可证

通过贡献，你同意你的贡献将在 MIT 许可证下授权。

## 致谢

感谢所有贡献者！你们的贡献让这个项目变得更好。
