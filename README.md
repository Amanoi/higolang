# HiGolang

Go 语言最新动态资讯聚合博客 —— 自动聚合 Go 语言相关的新闻、博客、版本发布等内容，为 Go 开发者提供一站式的资讯阅读体验。

## 功能特性

- **资讯聚合**：自动从多个 Go 语言相关的 RSS 源抓取最新文章
- **定时抓取**：内置调度器，支持自定义抓取间隔，自动保持内容更新
- **分类管理**：支持按主题/来源对文章进行分类和筛选
- **全文搜索**：支持关键词搜索，快速定位感兴趣的内容
- **管理后台**：提供完整的管理界面，支持内容管理、源管理和用户管理
- **响应式设计**：前端采用现代化 UI，支持桌面和移动端访问
- **用户认证**：基于 JWT 的身份认证，支持管理员和普通用户角色
- **SQLite 存储**：轻量级数据库，无需额外部署数据库服务

## 技术栈

### 后端
- **Go** + **Gin**：高性能 Web 框架
- **GORM**：Go ORM 库
- **SQLite**：轻量级嵌入式数据库（通过 CGO 驱动）
- **robfig/cron**：定时任务调度
- **golang-jwt**：JWT 认证

### 前端
- **Vue 3**：渐进式 JavaScript 框架
- **Vite**：下一代前端构建工具
- **Axios**：HTTP 客户端
- **Vue Router**：路由管理
- **Pinia**：状态管理

### 部署
- **Docker** + **Docker Compose**：容器化部署
- **Nginx**：前端静态资源服务与反向代理

## 快速开始

### Docker 部署（推荐）

确保已安装 [Docker](https://docs.docker.com/get-docker/) 和 [Docker Compose](https://docs.docker.com/compose/install/)。

```bash
# 克隆项目
git clone https://github.com/yourusername/higolang.git
cd higolang

# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

服务启动后：
- 前端页面：http://localhost
- 后端 API：http://localhost:8080

### 域名与 HTTPS 配置

项目支持通过环境变量配置域名和 SSL 证书，方便部署到不同环境。

**1. 配置域名**

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env，修改为你的域名
vi .env
```

`.env` 文件内容：

```bash
# 站点域名
DOMAIN=free.higolang.vip

# SSL 证书路径（容器内路径，留空则仅启用 HTTP）
SSL_CERT=/etc/nginx/ssl/fullchain.pem
SSL_KEY=/etc/nginx/ssl/privkey.pem
```

**2. 配置 SSL 证书**

将证书文件放入项目根目录的 `ssl/` 目录：

```bash
mkdir -p ssl
# 将你的证书文件放入 ssl/ 目录
cp /path/to/fullchain.pem ssl/
cp /path/to/privkey.pem ssl/
```

**3. 使用 Let's Encrypt 免费证书**

```bash
# 先以 HTTP 模式启动（.env 中 SSL_CERT 留空）
docker-compose up -d

# 安装 certbot 并申请证书
apt install certbot
certbot certonly --webroot -w certbot-webroot -d free.higolang.vip

# 将证书复制到 ssl/ 目录
cp /etc/letsencrypt/live/free.higolang.vip/fullchain.pem ssl/
cp /etc/letsencrypt/live/free.higolang.vip/privkey.pem ssl/

# 修改 .env 启用 SSL
sed -i 's|SSL_CERT=.*|SSL_CERT=/etc/nginx/ssl/fullchain.pem|' .env
sed -i 's|SSL_KEY=.*|SSL_KEY=/etc/nginx/ssl/privkey.pem|' .env

# 重启服务
docker-compose restart web
```

**4. 更换域名**

只需修改 `.env` 中的 `DOMAIN` 值并重启 web 容器：

```bash
# 修改域名
sed -i 's|DOMAIN=.*|DOMAIN=new.higolang.vip|' .env

# 重启生效
docker-compose restart web
```

### 使用 Make 命令

```bash
# 构建镜像
make build

# 启动服务
make up

# 查看日志
make logs

# 停止服务
make down

# 清理所有数据
make clean
```

## 本地开发

### 后端

确保已安装 [Go 1.24+](https://go.dev/dl/)。

```bash
cd server

# 安装依赖
go mod download

# 运行服务（默认端口 8080）
go run main.go
```

或使用 Make：

```bash
make dev-server
```

后端启动后，API 地址为 `http://localhost:8080`。

### 前端

确保已安装 [Node.js 18+](https://nodejs.org/)。

```bash
cd web

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

或使用 Make：

```bash
make dev-web
```

前端开发服务器默认运行在 `http://localhost:5173`。

## 默认管理员账号

首次启动时，系统会自动创建默认管理员账号（需确保 `config.yaml` 中 `seed.enabled` 为 `true`）：

- **用户名**：`admin`
- **密码**：`admin123`

> 请在首次登录后及时修改默认密码。生产环境部署前，请务必修改 `config.yaml` 中的 JWT 密钥。

## 默认数据源

系统预置以下 Go 语言相关的 RSS 数据源：

| 名称 | 地址 |
|------|------|
| Go 官方博客 | https://go.dev/blog/feed.atom |
| Go 版本发布 | https://go.dev/dl/?mode=json |
| Go Weekly | https://go.dev/weekly?format=rss |
| Golang News | https://golangnews.com/index.xml |
| Dave Cheney | https://dave.cheney.net/feed/xml |
| Go Time Podcast | https://changelog.com/gotime/feed |

可在管理后台中添加、编辑或删除数据源。

## 项目结构

```
higolang/
├── docker-compose.yml      # Docker Compose 配置
├── Makefile                # 常用命令快捷方式
├── README.md               # 项目说明文档
├── .gitignore              # Git 忽略规则
├── server/                 # 后端服务
│   ├── Dockerfile          # 后端 Docker 镜像
│   ├── config.yaml         # 应用配置文件
│   ├── main.go             # 入口文件
│   ├── go.mod              # Go 模块定义
│   ├── go.sum              # 依赖校验文件
│   ├── data/               # 数据存储目录
│   │   └── .gitkeep
│   ├── model/              # 数据模型
│   ├── handler/            # 请求处理器
│   ├── router/             # 路由定义
│   ├── scheduler/          # 定时任务调度
│   └── middleware/         # 中间件（JWT 认证等）
├── web/                    # 前端应用
│   ├── Dockerfile          # 前端 Docker 镜像
│   ├── package.json        # 依赖定义
│   ├── vite.config.js      # Vite 配置
│   ├── index.html          # 入口 HTML
│   └── src/                # 源代码
│       ├── main.js         # 应用入口
│       ├── App.vue         # 根组件
│       ├── views/          # 页面组件
│       ├── components/     # 通用组件
│       ├── router/         # 前端路由
│       ├── stores/         # 状态管理
│       └── api/            # API 请求封装
└── data/                   # 持久化数据目录（Docker 挂载）
```

## API 文档

所有接口前缀为 `/api/v1`，管理接口前缀为 `/api/v1/admin`。

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/articles` | 文章列表（支持 `page`、`page_size`、`category`、`tag`、`search` 参数） |
| GET | `/api/v1/articles/:slug` | 文章详情 |
| GET | `/api/v1/categories` | 分类列表 |
| GET | `/api/v1/tags` | 标签列表 |
| GET | `/api/v1/go-version` | 最新 Go 版本号 |
| GET | `/api/v1/settings/public` | 站点公开配置 |

### 管理接口（需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/admin/login` | 管理员登录 |
| GET | `/api/v1/admin/dashboard` | 仪表盘统计 |
| GET/POST | `/api/v1/admin/articles` | 文章列表 / 创建文章 |
| GET/PUT/DELETE | `/api/v1/admin/articles/:id` | 查看 / 编辑 / 删除文章 |
| GET/POST | `/api/v1/admin/categories` | 分类列表 / 创建分类 |
| PUT/DELETE | `/api/v1/admin/categories/:id` | 编辑 / 删除分类 |
| GET/POST | `/api/v1/admin/tags` | 标签列表 / 创建标签 |
| PUT/DELETE | `/api/v1/admin/tags/:id` | 编辑 / 删除标签 |
| GET/POST | `/api/v1/admin/feed-sources` | 数据源列表 / 添加数据源 |
| PUT/DELETE | `/api/v1/admin/feed-sources/:id` | 编辑 / 删除数据源 |
| POST | `/api/v1/admin/feed-sources/:id/fetch` | 手动触发抓取 |
| GET | `/api/v1/admin/fetch-logs` | 抓取日志 |
| GET/PUT | `/api/v1/admin/settings` | 查看 / 更新站点设置 |

> 管理接口需在请求头中携带 `Authorization: Bearer <token>`。

## 配置说明

`server/config.yaml` 主要配置项：

```yaml
server:
  port: 8080          # 服务端口
  mode: debug         # 运行模式：debug / release

database:
  path: ./data/higolang.db  # SQLite 数据库文件路径

jwt:
  secret: your-secret-key    # JWT 签名密钥（生产环境请修改）
  expire: 24                 # Token 有效期（小时）

scheduler:
  enabled: true              # 是否启用定时抓取
  default_interval: 30       # 默认抓取间隔（分钟）

seed:
  enabled: true              # 是否初始化默认数据（管理员账号 + 数据源）
```

## 许可证

本项目基于 [MIT 许可证](LICENSE) 开源。
