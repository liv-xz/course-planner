# 后端测试说明

## 目录
```
backend/test/
├── setup_test.go        # TestMain：建测试库、迁移、种子、清理
├── jwt_test.go          # JWT 单元测试（不需要数据库）
├── auth_api_test.go     # 登录 / profile 接口测试（需要数据库）
├── repository_test.go   # Repo CRUD 测试（需要数据库，自动清理）
└── README.md
```

## 运行
```bash
cd backend
go test ./test/ -v
```

只要单元测试（JWT）：
```bash
go test ./test/ -run TestGenerateAndParseToken -v
```

## 测试库

- 使用**独立的测试库** `course_planner_test`，**不会碰开发库 `course_planner`**。
- `TestMain` 每次运行会：
  1. `DROP DATABASE IF EXISTS course_planner_test`
  2. `CREATE DATABASE course_planner_test`
  3. AutoMigrate 建表
  4. 插入种子（角色 student/admin + 用户 admin/123456）
  5. 跑测试
  6. 结束后 `DROP DATABASE course_planner_test`（自动清理）

> 因此每次 `go test` 都是干净环境，跑完测试库自动删除。

## 依赖
- 需要本机 MySQL 已启动（`MySQL80` 服务）。
- 测试读取 `backend/configs/config.yaml` 的连接信息（host/user/password）。

## 测试清单
| 测试 | 说明 |
|---|---|
| `TestGenerateAndParseToken` | 生成/解析 token，claims 一致 |
| `TestParseToken_Invalid` | 非法 token 报错 |
| `TestParseToken_WrongSecret` | 密钥不一致报错 |
| `TestParseToken_Expired` | 过期 token 报错 |
| `TestLogin_Success` | admin/123456 登录成功，返回 token |
| `TestLogin_WrongPassword` | 错误密码 401 |
| `TestProfile_WithToken` | 带 token 访问 profile 成功 |
| `TestProfile_NoToken` | 不带 token 401 |
| `TestRoleRepo_CRUD` | 角色增删查改 |
| `TestUserRepo_SoftDelete` | 用户逻辑删除 |
| `TestPrerequisiteRepo` | 先修关系查询 |
