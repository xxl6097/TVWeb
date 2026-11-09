# TVWeb 升级功能说明

## 概述

TVWeb 集成了自动升级检测功能，可以自动检查 GitHub Releases 并提供更新通知。

## API 接口

### 1. 检查更新

**接口**: `GET http://localhost:8091/api/upgrade/check`

**响应示例**:

```json
{
  "hasUpdate": true,
  "currentVersion": "v0.0.5",
  "latestVersion": "v0.0.6",
  "releaseNote": {
    "version": "v0.0.6",
    "notes": "Release notes in markdown format",
    "publishedAt": "2024-11-09 21:30:00",
    "pageUrl": "https://github.com/xxl6097/TVWeb/releases/tag/v0.0.6",
    "fullUrl": "https://github.com/xxl6097/TVWeb/releases/download/v0.0.6/TVWeb_v0.0.6_darwin_arm64",
    "patchUrl": "",
    "size": 12345678
  }
}
```

### 2. 查看版本信息

**接口**: `GET http://localhost:8091/api/version`

**响应示例**:

```json
{
  "appName": "TVWeb",
  "appVersion": "v0.0.5",
  "buildTime": "2024-11-09 21:25:30",
  "gitCommit": "2af1001",
  "gitBranch": "main",
  "osType": "darwin",
  "arch": "arm64",
  "githubUser": "xxl6097",
  "githubRepo": "TVWeb"
}
```

## 使用方法

### 命令行测试

```bash
# 检查更新
curl http://localhost:8091/api/upgrade/check

# 查看版本信息
curl http://localhost:8091/api/version
```

### JavaScript/前端集成

```javascript
// 检查更新
async function checkForUpdates() {
  try {
    const response = await fetch('http://localhost:8091/api/upgrade/check');
    const data = await response.json();
    
    if (data.hasUpdate) {
      console.log(`发现新版本: ${data.latestVersion}`);
      console.log(`当前版本: ${data.currentVersion}`);
      console.log(`下载地址: ${data.releaseNote.pageUrl}`);
      
      // 显示更新提示
      showUpdateNotification(data);
    } else {
      console.log('已是最新版本');
    }
  } catch (error) {
    console.error('检查更新失败:', error);
  }
}

// 获取版本信息
async function getVersion() {
  const response = await fetch('http://localhost:8091/api/version');
  const version = await response.json();
  console.log(`当前版本: ${version.appVersion}`);
  return version;
}
```

### HTML 示例界面

```html
<!DOCTYPE html>
<html>
<head>
  <title>TVWeb 升级检查</title>
  <style>
    .update-card {
      padding: 20px;
      border: 1px solid #ddd;
      border-radius: 8px;
      margin: 20px;
    }
    .btn {
      padding: 10px 20px;
      background: #007bff;
      color: white;
      border: none;
      border-radius: 4px;
      cursor: pointer;
    }
  </style>
</head>
<body>
  <div class="update-card">
    <h2>TVWeb 版本管理</h2>
    <div id="version-info"></div>
    <button class="btn" onclick="checkUpdate()">检查更新</button>
    <div id="update-info" style="margin-top: 20px;"></div>
  </div>

  <script>
    async function loadVersion() {
      const res = await fetch('http://localhost:8091/api/version');
      const data = await res.json();
      document.getElementById('version-info').innerHTML = `
        <p>当前版本: <strong>${data.appVersion}</strong></p>
        <p>构建时间: ${data.buildTime}</p>
        <p>平台: ${data.osType}/${data.arch}</p>
      `;
    }

    async function checkUpdate() {
      const res = await fetch('http://localhost:8091/api/upgrade/check');
      const data = await res.json();
      const infoDiv = document.getElementById('update-info');
      
      if (data.hasUpdate) {
        infoDiv.innerHTML = `
          <div style="background: #d4edda; padding: 15px; border-radius: 4px;">
            <h3>🎉 发现新版本！</h3>
            <p>最新版本: <strong>${data.latestVersion}</strong></p>
            <p>发布时间: ${data.releaseNote.publishedAt}</p>
            <a href="${data.releaseNote.pageUrl}" target="_blank" class="btn">
              查看更新说明
            </a>
          </div>
        `;
      } else {
        infoDiv.innerHTML = `
          <div style="background: #d1ecf1; padding: 15px; border-radius: 4px;">
            <p>✅ 已是最新版本</p>
          </div>
        `;
      }
    }

    // 页面加载时显示版本信息
    loadVersion();
  </script>
</body>
</html>
```

## 端口说明

- **主服务端口**: 默认 `:8088` (可配置)
- **升级 API 端口**: `:8091` (固定)

## 自动检查

服务启动时会自动检查一次更新，如果发现新版本会在日志中输出。

## 工作原理

1. 服务启动时初始化 Upgrader
2. 通过 GitHub API 检查最新 Release
3. 比较当前版本和最新版本
4. 提供 HTTP API 供前端调用
5. 支持手动触发升级检查

## 注意事项

- 需要网络连接访问 GitHub API
- 自动升级功能依赖 go-service 框架的 Upgrade 方法
- 建议在前端定期调用检查接口（如每小时检查一次）
