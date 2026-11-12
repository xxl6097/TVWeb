# 发布流程自动化脚本

请按以下顺序执行任务：

## 步骤 1：更新上游依赖
检查并更新 github.com/xxl6097/TVBoxWeb 依赖到最新版本：
- 运行 `go get -u github.com/xxl6097/TVBoxWeb@latest`
- 运行 `go mod tidy`
- 确认依赖更新成功

## 步骤 2：更新 Release.md
修改 Release.md 文件中的以下两个部分（保持"github加速"部分不变）：
- **新特性**：[根据代码变化识别新增了哪些新特性功能]
- **修复内容**：[根据代码变化识别修复的bug列表]

注意：直接从上游仓库 commit/release 记录中提取本次变更内容，不需要询问用户。

## 步骤 3：提交代码并发布
- 使用 `git add .` 暂存所有变更
- 创建 commit，message 格式：`Release vX.X.X: [提炼的新功能和修复的bug]`
- 创建新 tag，格式：`vX.X.X`（根据语义化版本规则递增，但是要高于当前release版本号）
- 推送代码和 tag：`git push && git push --tags`

注意：发布的 release notes 会自动使用 Release.md 的内容。
