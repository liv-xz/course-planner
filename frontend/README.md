# 课程规划系统前端

上海建桥学院毕业设计项目的前端管理平台，基于 Vue 3 + Vite + Element Plus。

## 已实现功能

- 登录与角色区分：管理员 / 学生
- 管理员工作台
- 培养方案增删改查
- 课程类别树增删改查
- 课程库增删改查
- 课程-类别归属（叶子类别、同一方案单归属）
- 先修关系维护
- 开课学期维护
- 专业管理
- 学期字典管理
- 学生管理（新增、编辑、删除、搜索）
- 学生修读记录增删改查
- 管理员查看学生学分审核
- 学生端：我的修读记录、学分审核、选课推荐
- 智能体：对话、SSE、工具调用提示、会话管理
- 管理员：学分缺口统计、安全拦截日志
- 修读记录 CSV/Excel 导入

完整项目文档见 [../项目文档/README.md](../项目文档/README.md)。

## 启动

先启动后端：

```bash
cd backend
go run ./cmd/server
```

再启动前端：

```bash
cd frontend
npm install
npm run dev
```

访问 `http://localhost:5173`。

Vite 已配置 `/api` 代理到 `http://127.0.0.1:8080`，所以不需要额外配置 CORS。

## 配置

后端配置在 `backend/configs/config.local.yaml`（已加入 .gitignore）：

- `database`：数据库连接
- `jwt.secret`：JWT 密钥
- `llm.api_key`：大模型 API Key，填入后把 `llm.enabled` 改为 `true`
- `security.patterns`：注入拦截关键词
- `import`：导入文件大小/行数/后缀

不填 LLM Key 也能运行，智能体会自动使用规则兜底。

## 演示账号

| 角色 | 账号 | 密码 |
|---|---|---|
| 教务管理员 | `admin` | `123456` |
| 学生 | `2323500` | `123456` |

## 构建

```bash
npm run build
```

产物在 `dist/`，可以由任意静态服务器托管；生产环境需要把 `/api` 反向代理到后端。
