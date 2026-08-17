# 记账系统（Bookkeeping）

一个纯 Go 标准库实现的个人记账 REST API 服务，采用标准 Go 工程目录结构，内存存储，零第三方依赖。

> 金额统一使用 int64 的「分」为单位，避免浮点精度问题。

## 目录结构

```
origin/
├── cmd/server/          # 程序入口
├── internal/
│   ├── app/             # 依赖装配
│   ├── config/          # 配置加载
│   ├── model/           # 领域模型与校验
│   ├── store/           # 数据访问接口 + 内存实现
│   ├── service/         # 业务逻辑层
│   └── handler/         # HTTP 处理器层
└── pkg/
    ├── httpx/           # HTTP 响应工具
    ├── idgen/           # ID 生成
    └── logger/          # 分级日志
```

## 运行

```bash
go run ./cmd/server
PORT=8081 go run ./cmd/server
```

默认监听 `:8080`。

## 测试

```bash
go test ./...
```

## API 接口

### 账户

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/accounts` | 创建账户 |
| GET | `/api/accounts` | 账户列表 |
| GET | `/api/accounts/{id}` | 账户详情 |
| PUT | `/api/accounts/{id}` | 更新账户 |
| DELETE | `/api/accounts/{id}` | 删除账户 |

### 分类

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/categories` | 创建分类 `{"name":"餐饮","type":"expense"}` |
| GET | `/api/categories?type=expense` | 分类列表 |
| GET | `/api/categories/{id}` | 分类详情 |
| PUT | `/api/categories/{id}` | 更新分类 |
| DELETE | `/api/categories/{id}` | 删除分类 |

### 流水

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/transactions` | 记一笔（自动更新余额） |
| GET | `/api/transactions?account_id=&category_id=&type=&from=&to=&page=&size=` | 流水列表 |
| GET | `/api/transactions/{id}` | 流水详情 |
| DELETE | `/api/transactions/{id}` | 删除流水（回滚余额） |

### 预算

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/budgets` | 创建预算 |
| GET | `/api/budgets` | 预算列表 |
| GET | `/api/budgets/{id}` | 预算详情 |
| PATCH | `/api/budgets/{id}` | 更新预算金额 |
| DELETE | `/api/budgets/{id}` | 删除预算 |
| GET | `/api/budgets/check` | 检查当月预算执行情况 |

### 统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/stats/overview` | 账户总览 |
| GET | `/api/stats/monthly?year=&month=` | 月度收支报表 |
