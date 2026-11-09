# TVWeb Web 界面扩展

本目录包含 TVWeb 的 Web 界面扩展组件。

## 文件说明

### 1. `upgrade.html` - 独立的升级检测页面

**访问地址**: 将此文件放到 TVBoxWeb 的静态资源目录，或者直接用浏览器打开

**功能特性**:
- ✅ 显示当前版本信息
- 🔄 一键检查更新
- 🎉 发现新版本时显示详细信息
- 📥 直接跳转到 GitHub Release 页面下载

**使用方法**:

方法一：直接在浏览器打开
```bash
# 确保 TVWeb 服务正在运行
open web/upgrade.html
```

方法二：集成到 TVBoxWeb
```bash
# 将文件复制到 TVBoxWeb 的静态目录
cp web/upgrade.html <tvboxweb-static-dir>/upgrade.html
# 然后访问 http://localhost:8088/upgrade.html
```

### 2. `UpgradePanel.vue` - Vue 3 组件

**适用场景**: 如果你有 TVBoxWeb 的源码并想直接集成到 Vue 应用中

**集成步骤**:

1. 将 `UpgradePanel.vue` 复制到 TVBoxWeb 项目的组件目录
```bash
cp web/UpgradePanel.vue /path/to/TVBoxWeb/web/src/components/
```

2. 在设置页面中引入组件
```vue
<script setup>
import UpgradePanel from '@/components/UpgradePanel.vue'
</script>

<template>
  <div class="settings-page">
    <!-- 其他设置选项 -->
    
    <!-- 升级检测面板 -->
    <section class="settings-section">
      <h2>系统更新</h2>
      <UpgradePanel />
    </section>
  </div>
</template>
```

3. 或者在路由中添加独立页面
```javascript
// router/index.js
{
  path: '/upgrade',
  name: 'upgrade',
  component: () => import('@/components/UpgradePanel.vue')
}
```

## API 端点

升级检测依赖以下 API：

- `GET http://localhost:8091/api/version` - 获取版本信息
- `GET http://localhost:8091/api/upgrade/check` - 检查更新

确保 TVWeb 服务已启动，这些 API 会自动可用。

## 快速开始

1. 启动 TVWeb 服务
```bash
./TVWeb start
# 或者使用已编译的二进制文件
```

2. 用浏览器打开升级页面
```bash
open web/upgrade.html
```

3. 点击"检查更新"按钮

## 截图

升级页面包含：
- 📊 版本信息卡片（当前版本、构建时间、平台架构）
- 🔄 检查更新按钮
- 🎉 新版本提示（如果有更新）
- ✅ 最新版本提示（已是最新）
- ⚠️ 错误提示（检查失败时）

## 自定义

### 修改 API 地址

如果你的升级 API 端口不是 8091，修改文件中的 API 地址：

```javascript
// upgrade.html 中
const res = await fetch('http://localhost:YOUR_PORT/api/version');

// UpgradePanel.vue 中
const res = await fetch('http://localhost:YOUR_PORT/api/version');
```

### 修改样式

两个文件都使用了 CSS 变量，你可以轻松自定义颜色：

```css
:root {
  --primary: #667eea;
  --primary-hover: #764ba2;
  --success: #00c853;
  --error: #c62828;
}
```

## 注意事项

- 升级 API 运行在独立端口 `:8091`
- 主 TVBoxWeb 服务默认运行在 `:8088`
- 如果出现 CORS 错误，确保两个服务都在本地运行
- 检查更新需要网络连接访问 GitHub API

## 问题排查

**Q: 显示"加载中..."不变**
A: 检查 TVWeb 服务是否启动，API 端口是否正确

**Q: 检查更新失败**
A: 确认网络连接，检查是否能访问 GitHub

**Q: 版本信息不显示**
A: 确保使用最新版本的 TVWeb，并使用构建脚本编译

## 更多信息

详细的 API 文档和集成说明请参考 `README_UPGRADE.md`
