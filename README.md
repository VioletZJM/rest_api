![ci](https://github.com/VioletZJM/rest_api/actions/workflows/ci.yml/badge.svg)

# REST API

基于 Gin + SQLite 的事件（Event）管理 API：用户注册 / 登录（JWT）、事件 CRUD、活动报名与参会人查询。

## 技术栈

Go 1.26.5 · Gin · SQLite（`mattn/go-sqlite3`，**需要 CGO**）· JWT（HS256）· bcrypt

## 快速开始

### 本地运行

前置：Go 1.26.5+、C 编译器（go-sqlite3 依赖 CGO）。

```bash
export JWT_SECRET=your-secret   # 必填，未设置会直接启动失败
make run                        # 等价于 go run .，监听 :8090
```

首次启动会自动创建 `api.db`，并建好 `users` / `events` / `registrations` 三张表。

跑测试：

```bash
make test     # go test -count=1 ./...
make cover    # 生成 cover.out 并输出总覆盖率
```

### Docker

```bash
JWT_SECRET=your-secret make up   # 构建并后台启动，数据写入 api-data 卷
make logs                        # 跟踪日志
make down                        # 停止并移除容器
```

> `docker-compose.yml` 在未提供 `JWT_SECRET` 时会回退为 `helloworld`，正式环境务必覆盖。
>
> 其他常用命令：`make build`（编译到 `bin/server`）、`make tidy`（整理依赖）。

## 配置

| 环境变量 | 默认值 | 说明 |
|---|---|---|
| `JWT_SECRET` | 无（必填） | JWT 签名密钥 |
| `PORT` | `8090` | 监听端口 |
| `DB_PATH` | `api.db` | SQLite 文件路径 |

## 接口

Base URL：`http://localhost:8090`。标注 🔒 的接口需携带请求头 `Authorization: Bearer <token>`。

| 方法 | 路径 | 鉴权 | 说明 | 成功状态 |
|---|---|---|---|---|
| POST | `/signup` | | 注册，body：`email`、`password` | 201 |
| POST | `/login` | | 登录，返回 `message`、`token` | 200 |
| GET | `/events` | | 事件列表，支持分页 | 200 |
| GET | `/events/:id` | | 事件详情 | 200 |
| POST | `/events` | 🔒 | 创建事件，body：`name`、`description`、`location`、`dateTime` | 201 |
| PUT | `/events/:id` | 🔒 | 更新事件（仅创建者） | 200 |
| DELETE | `/events/:id` | 🔒 | 删除事件（仅创建者） | 200 |
| POST | `/events/:id/register` | 🔒 | 报名活动 | 201 |
| DELETE | `/events/:id/register` | 🔒 | 取消报名 | 200 |
| GET | `/events/:id/attendees` | 🔒 | 参会人列表 | 200 |
| GET | `/registrations/me` | 🔒 | 我报名的事件 | 200 |

最小调用示例：

```bash
curl -X POST localhost:8090/signup -H 'Content-Type: application/json' \
  -d '{"email":"a@example.com","password":"123456"}'

TOKEN=$(curl -s -X POST localhost:8090/login -H 'Content-Type: application/json' \
  -d '{"email":"a@example.com","password":"123456"}' | sed -E 's/.*"token":"([^"]+)".*/\1/')

curl -X POST localhost:8090/events -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Meetup","description":"desc","location":"online","dateTime":"2026-09-21T12:00:00Z"}'
```

### 分页

`GET /events?page=1&limit=20`

- `page` 默认 1，非法值回退为 1；`limit` 默认 20，非法或小于 1 回退为 20，上限 100。
- 响应结构：`{"total":25,"page":1,"limit":20,"Items":[...]}`（注意 `Items` 为大写开头，其余为小写）。
- 超出末页返回空 `Items`，不返回 404。

### 错误格式

所有错误统一为：

```json
{ "error": { "code": "NOT_FOUND", "message": "resource not found" } }
```

| 状态码 | code | 场景 |
|---|---|---|
| 400 | `INVALID_REQUEST_BODY` / `INVALID_EVENT_ID` | 请求体不合法 / 路径 id 非数字 |
| 401 | `unauthorized` / `INVALID_CREDENTIALS` | 缺少或无效 token / 邮箱或密码错误 |
| 403 | `forbidden` | 非事件创建者执行修改或删除 |
| 404 | `NOT_FOUND` / `REGISTRATION_NOT_FOUND` | 事件不存在 / 未报名却取消报名 |
| 409 | `ALREADY_REGISTERED` | 重复报名同一活动 |
| 500 | `INTERNAL_SERVER_ERROR` | 其他服务端错误（详情只写日志） |

## 项目结构

```
.
├── main.go                 # 入口：读配置、加载密钥、初始化 DB、启动与优雅关闭
├── database/db.go          # SQLite 连接、连接池、建表与唯一索引
├── models/
│   ├── user.go             # 用户模型：注册、凭据校验
│   └── event.go            # 事件模型：CRUD、报名、参会人及哨兵错误
├── routes/
│   ├── routes.go           # 路由注册（公开 / 需鉴权分组）
│   ├── users.go            # /signup、/login
│   ├── events.go           # 事件 CRUD、分页列表、参会人
│   ├── register.go         # 报名 / 取消报名
│   ├── parse.go            # 分页参数解析与兜底
│   ├── errors.go           # 哨兵错误 → HTTP 状态码 / 错误码映射
│   └── *_test.go           # handler 层端到端测试（临时 SQLite）
├── middlewares/auth.go     # Bearer Token 鉴权，写入 userId
├── utils/
│   ├── config.go           # 环境变量读取与校验
│   ├── jwt.go              # 签发 / 校验 token
│   ├── hash.go             # bcrypt 哈希与比对
│   └── respond.go          # 统一错误响应
├── api_test/*.http         # 手工调试脚本（VS Code REST Client）
├── Dockerfile              # 多阶段构建，非 root 运行
├── docker-compose.yml      # 单服务 + 命名卷
├── Makefile                # 常用命令入口
└── .github/workflows/ci.yml
```

## 项目说明

- **来源**：本项目基于开源 Go REST API 教学项目搭建，在其基础上补充了分页、统一错误响应、JWT 鉴权、测试与容器化。
- **JSON 字段名**：`Event` / `User` / `Attendee` 未加 `json` tag，因此响应中的字段名首字母大写（`ID`、`Name`、`Location`、`UserID` 等）；请求体反序列化不区分大小写，`dateTime` 可直接使用。
- **时间格式**：`dateTime` 使用 RFC3339，如 `2026-09-21T12:00:00Z`。
- **鉴权**：登录返回的 JWT 有效期 2 小时；middleware 校验后将 `userId` 放入上下文供 handler 使用。
- **权限**：仅事件创建者可 `PUT` / `DELETE` 该事件；`GET :id/attendees` 只要求登录。
- **报名唯一性**：`registrations(event_id, user_id)` 建有唯一索引，重复报名返回 409；取消未报名记录返回 404。
- **数据层约定**：`models` 同时承担 SQL 访问职责，并定义 `ErrNotFound` 等哨兵错误，由 `routes/errors.go` 统一映射为 HTTP 响应。
- **测试**：每个用例使用独立临时 SQLite 文件，需设置 `JWT_SECRET`（CI 中为 `ci-test-secret`）；CI 执行 build / vet / test，并把覆盖率写入 Job Summary、上传 `cover.out`。
- **运维**：服务收到 `SIGINT` / `SIGTERM` 后最多等待 5 秒完成优雅关闭；容器以非 root 用户运行，数据库落在 `/data` 卷。
