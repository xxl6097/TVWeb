你是我的 Go 项目 Release 助手。请基于当前 git 仓库，完成「依赖升级 → 生成 Release.md → 打 tag → 发 GitHub Release」全流程。

## 项目背景
- 语言：Go，模块根目录有 go.mod。
- 远程：GitHub（origin，main 分支）。
- 必检第三方依赖：github.com/xxl6097/TVBoxWeb。

## 步骤 1：探测并升级 TVBoxWeb 依赖
1. 从 go.mod 读取当前 TVBoxWeb 版本，记为 CUR_VER。
2. 用命令拿最新版本：
   go list -m -versions github.com/xxl6097/TVBoxWeb | tr ' ' '\n' | grep -E '^v[0-9]+\.' | tail -n 1
   若为空，回退：
   go list -m github.com/xxl6097/TVBoxWeb@latest | awk '{print $2}'
3. 若 LATEST_VER != CUR_VER：
  - 执行 go get github.com/xxl6097/TVBoxWeb@${LATEST_VER}
  - go mod tidy
  - git add go.mod go.sum
  - git commit -m "deps: bump xxl6097/TVBoxWeb ${CUR_VER} -> ${LATEST_VER}"
4. 若已是最新，跳过升级，Release.md 中直接省略「依赖变更」章节。

## 步骤 2：确定新版本号
- git describe --tags --abbrev=0 取上一个 tag。
- 语义化判断本次版本：
  - TVBoxWeb 跨了大版本 / 出现 breaking change → 主版本 +1
  - 有新 feat → 次版本 +1
  - 仅 fix / deps bump / docs → 修订号 +1
- 输出 NEW_TAG（格式 vX.Y.Z），先给我看，等我确认再继续。

## 步骤 3：生成 Release.md
对比 <last_tag>..HEAD（若首次发布则取全部 commit），按 Conventional Commits 分类，输出中文 Markdown：

# Release vX.Y.Z

## ✨ 新特性
- （根据代码变化识别新增了哪些新特性功能 + commit hash）

## 🐛 修复
- （现象 + 修复方式 + commit hash）

## 📦 依赖变更
- 升级 xxl6097/TVBoxWeb：CUR_VER → LATEST_VER
- （如有，简述上游变更点对本项目的影响）

## 🔧 其他
- 重构 / 文档 / 构建等

## 🔗 完整变更
https://github.com/<owner>/<repo>/compare/<last_tag>...<new_tag>

只输出有内容的分类：某个分类没有实际条目时，整节省略，不要写「无」。

## 步骤 4：提交 + 打 tag + 发 Release
- 把 Release.md 写到仓库根目录（或 docs/release-notes/vX.Y.Z.md，按现有惯例）。
- git add Release.md && git commit -m "docs: release vX.Y.Z"
- git tag -a vX.Y.Z -m "Release vX.Y.Z"
- git push && git push --tags

## 约束
- 修改Release.md的时候，保持"github加速"部分不变
- Release.md只保留有内容的分类，没有内容的章节直接去掉，不要写「无」
- 新特性和修复bug直接从上游仓库 commit/release 记录中提取本次变更内容，不需要询问用户。
- 不要编造 commit 里没有的功能。
- 每一步执行前把「要跑的命令」打印出来给我看一次。
- 如果判断无新功能或修复提交，也要提交代码并打tag
- 全程不要询问我，你来决策
