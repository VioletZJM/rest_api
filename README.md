![ci](https://github.com/VioletZJM/rest_api/actions/workflows/ci.yml/badge.svg)

# REST API
本项目基于开源Go REST API教学项目搭建
基于 Gin + SQLite 的事件管理 API：用户注册/登录（JWT）、事件 CRUD、活动报名。

## 技术栈

Go 1.26+、Gin、SQLite（`github.com/mattn/go-sqlite3`，**需要 CGO**）、JWT、bcrypt。

## 快速开始

```bash
export JWT_SECRET=your-secret   # 必填，未设置会启动失败
go run .                        # 监听 http://localhost:8090
```

数据库文件 `api.db` 自动创建，表结构在启动时初始化（users / events / registrations）。

测试：

```bash
JWT_SECRET=ci-test-secret go test ./...
```

## 项目结构

```
.
├── main.go                 # 入口：加载密钥、初始化 DB、启动服务（:8090）
├── database/db.go          # SQLite 连接与建表
├── models/                 # 数据模型与 SQL 操作
│   ├── user.go             # 用户、密码校验
│   └── event.go            # 事件、报名、参会人
├── routes/                 # 路由与 handler
│   ├── routes.go           # 路由注册
│   ├── users.go            # 注册 / 登录
│   ├── events.go           # 事件 CRUD、报名查询
│   ├── register.go         # 报名 / 取消报名
│   └── auth_test.go
├── middlewares/auth.go     # JWT 鉴权中间件
├── utils/                  # JWT 签发校验、bcrypt 哈希
└── api_test/*.http         # 手工接口调试脚本
```

## 接口

基础地址：`http://localhost:8090`。带 🔒 的接口需在 `Authorization` 请求头中携带 `Bearer` + 空格 + token。

| 方法 | 路径 | 说明 | 成功状态 |
|---|---|---|---|
| POST | `/signup` | 注册，body: `{"email","password"}` | 201 |
| POST | `/login` | 登录，返回 `{"message","token"}` | 200 |
| GET | `/events` | 事件列表 | 200 |
| GET | `/events/:id` | 事件详情 | 200 |
| POST | 🔒 `/events` | 创建事件，body: `{"name","description","location","dateTime"}` | 201 |
| PUT | 🔒 `/events/:id` | 更新事件（仅创建者） | 200 |
| DELETE | 🔒 `/events/:id` | 删除事件（仅创建者） | 200 |
| POST | 🔒 `/events/:id/register` | 报名活动 | 201 |
| DELETE | 🔒 `/events/:id/register` | 取消报名 | 200 |
| GET | 🔒 `/events/:id/attendees` | 活动报名用户 `[{"ID","Email"}]` | 200 |
| GET | 🔒 `/registrations/me` | 我报名的事件 | 200 |

说明：

- 事件对象字段：`ID`、`Name`、`Description`、`Location`、`DateTime`、`UserID`。
- `dateTime` 使用 RFC3339，如 `2026-09-21T12:00:00Z`。
- 错误统一返回 `{"message": "..."}`：400 参数错误、401 未认证/非创建者、409 重复报名、500 服务端错误（含事件不存在）。
- Token 有效期 2 小时。
