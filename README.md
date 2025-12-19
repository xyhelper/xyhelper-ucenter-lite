# xyhelper项目模板

## 目录结构

- `.devcontainer/`: 开发容器配置文件
  - 包含开发环境数据库配置
  - 常用工具和脚本
  - 开发环境数据库存储

- `.vscode/`: VSCode配置文件
  - 调试配置
  - 插件推荐
  - 代码片段

- `backend/`: 后端代码目录
  - 启动命令：`cd backend && go run main.go`

- `frontend/`: 前端代码目录
  - 启动命令：`cd frontend && npm run dev`

## 开发指南

### 开发模式外网访问

#### 前端访问
```bash
cloudflared -url http://127.0.0.1:9000
```

#### 后端访问
```bash
cloudflared -url http://127.0.0.1:8001
```

### 常用命令

项目使用 makefile 简化常用命令：

| 命令 | 说明 |
|------|------|
| `make init` | 安装前后端依赖 |
| `make b` | 开发模式启动后端 |
| `make f` | 开发模式启动前端 |
| `make db` | 启动开发数据库 |

## 常见问题

### Git 配置
如果遇到 git 缺少用户名密码无法提交的问题，请运行以下命令：

```bash
git config --global user.email "you@example.com"
git config --global user.name "Your Name"
```