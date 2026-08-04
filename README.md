# iSCSI Web Panel

一个现代化的 iSCSI Target 管理面板，提供 Web UI 和 RESTful API，让 iSCSI 存储管理变得简单。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/go-1.21+-00ADD8.svg)
![Platform](https://img.shields.io/badge/platform-Linux%20x86_64%20%7C%20arm64-green.svg)

## 功能特性

- 🎯 **Target 管理** - 创建、删除、查看 iSCSI Target
- 💾 **LUN 管理** - 配置逻辑单元，支持多种后端存储
- 🔐 **CHAP 认证** - 完整的 initiator 认证和访问控制
- 📊 **实时监控** - CPU、内存、磁盘、网络、IOPS 指标
- 🚨 **告警系统** - 可配置的告警规则和通知
- 📸 **快照管理** - 创建和恢复 LUN 快照
- 👥 **用户管理** - 多用户角色权限控制
- 📝 **操作日志** - 完整的操作审计日志
- 🌙 **暗色模式** - 自动跟随系统主题
- 📱 **移动适配** - 响应式设计，支持手机访问

## 快速开始

### 一键部署

```bash
# 下载并运行
curl -fsSL https://raw.githubusercontent.com/Level6me/iscsi-web-panel/main/install.sh | sudo bash
```

或手动部署：

```bash
git clone https://github.com/Level6me/iscsi-web-panel.git
cd iscsi-web-panel
sudo ./install.sh
```

### 手动安装

```bash
# 1. 安装依赖
sudo apt-get update
sudo apt-get install -y tgt

# 2. 编译
go build -o iscsi-web-panel

# 3. 运行
sudo ./iscsi-web-panel
```

### Docker 部署

```bash
docker build -t iscsi-web-panel .
docker run -d --name iscsi-panel \
  --privileged \
  -p 3005:3005 \
  -v /sys/kernel/config:/sys/kernel/config \
  iscsi-web-panel
```

## 访问面板

部署完成后，访问：

```
http://<服务器IP>:3005
```

默认管理员账号：
- 用户名: `admin`
- 密码: `admin123`

⚠️ **重要**: 首次登录后请立即修改默认密码！

## 配置

环境变量配置：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `LISTEN_ADDR` | `:3005` | 监听地址 |
| `JWT_SECRET` | `iscsi-web-panel-secret-key` | JWT 密钥 |
| `DATA_DIR` | `./data` | 数据目录 |
| `DB_PATH` | `./data/iscsi.db` | 数据库路径 |

## API 文档

访问 `http://<服务器IP>:3005/api/v1/api-doc` 查看完整 API 文档。

### 主要接口

- `POST /api/v1/auth/login` - 登录
- `GET /api/v1/dashboard/overview` - 仪表盘概览
- `GET/POST /api/v1/targets` - Target 管理
- `GET/POST /api/v1/luns` - LUN 管理
- `GET/POST /api/v1/initiators` - Initiator 管理
- `GET /api/v1/monitor/metrics` - 监控指标
- `GET /api/v1/alerts` - 告警列表

## 系统要求

- Linux x86_64 / arm64
- tgt (SCSI Target Framework)
- systemd (用于服务管理)

## 技术栈

- **后端**: Go + Gin
- **前端**: 原生 HTML/CSS/JS (iOS 风格)
- **数据库**: SQLite
- **认证**: JWT
- **iSCSI**: tgt (tgtadm)

## 开发

```bash
# 克隆仓库
git clone https://github.com/Level6me/iscsi-web-panel.git
cd iscsi-web-panel

# 安装依赖
go mod download

# 运行开发服务器
go run main.go

# 重新生成前端
cd frontend
python3 build.py
```

## 截图

_待添加_

## License

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
