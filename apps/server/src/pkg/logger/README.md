# 日志系统 (Logging System)

基于 Winston 和 Morgan 的企业级日志系统，提供结构化日志记录、自动轮转、上下文管理等功能。

## 特性

- 🏗️ **结构化日志**: JSON 格式，便于分析和查询
- 🔄 **自动轮转**: 按日期和大小自动轮转日志文件
- 📊 **多级别日志**: 支持 error、warn、info、http、debug 等级别
- 🎯 **上下文管理**: 为不同模块创建带上下文的日志器
- 📡 **HTTP 日志**: 自动记录 HTTP 请求和响应
- 🆔 **请求追踪**: 自动生成请求 ID 用于链路追踪
- ⚡ **高性能**: 异步写入，不阻塞主线程
- 🎨 **开发友好**: 开发环境彩色输出，生产环境结构化存储

## 安装

```bash
pnpm add winston winston-daily-rotate-file morgan
pnpm add -D @types/morgan
```

## 快速开始

### 1. 基础日志使用

```typescript
import { logger, createContextLogger } from '@pkg/logger';

// 基础日志
logger.info('Application started');
logger.error('Something went wrong', { error: 'details' });

// 上下文日志
const userLogger = createContextLogger('UserService');
userLogger.info('User created', { userId: '123' });
```

### 2. Express 集成

```typescript
import express from 'express';
import { createMorganMiddleware, requestIdMiddleware } from '@pkg/logger/morgan';

const app = express();

// 添加请求 ID 和 HTTP 日志
app.use(requestIdMiddleware);
app.use(createMorganMiddleware());

app.get('/users', (req, res) => {
  const logger = createContextLogger('UserController');
  logger.info('Fetching users', { requestId: req.id });
  
  res.json({ users: [] });
});
```

## 日志级别

| 级别 | 数值 | 描述 |
|------|------|------|
| error | 0 | 错误信息 |
| warn | 1 | 警告信息 |
| info | 2 | 一般信息 |
| http | 3 | HTTP 请求 |
| verbose | 4 | 详细信息 |
| debug | 5 | 调试信息 |
| silly | 6 | 最详细 |

## 配置

通过环境变量配置：

```bash
# 日志级别
LOG_LEVEL=info

# 环境
NODE_ENV=production
```

## 日志文件

日志文件自动存储在 `logs/` 目录：

```
logs/
├── application-2023-12-01.log  # 所有日志
├── error-2023-12-01.log        # 错误日志
├── http-2023-12-01.log         # HTTP 请求日志
├── exceptions-2023-12-01.log   # 未捕获异常
└── rejections-2023-12-01.log   # 未处理的 Promise 拒绝
```

## API 文档

### Logger

#### 基础方法

```typescript
logger.error(message: string, meta?: Error | object): void
logger.warn(message: string, meta?: object): void
logger.info(message: string, meta?: object): void
logger.http(message: string, meta?: object): void
logger.debug(message: string, meta?: object): void
```

#### 示例

```typescript
// 简单日志
logger.info('User login successful');

// 带元数据的日志
logger.info('User login', {
  userId: '123',
  ip: '192.168.1.1',
  userAgent: 'Chrome/91.0'
});

// 错误日志
try {
  // some operation
} catch (error) {
  logger.error('Operation failed', error);
}
```

### ContextLogger

```typescript
const userLogger = createContextLogger('UserService');

// 所有日志会自动添加 [UserService] 前缀
userLogger.info('Creating user');  // 输出: [UserService] Creating user
userLogger.error('Validation failed', { field: 'email' });
```

### Morgan 中间件

#### createMorganMiddleware()

创建 HTTP 请求日志中间件：

```typescript
app.use(createMorganMiddleware());
```

#### requestIdMiddleware

添加请求 ID 和计时：

```typescript
app.use(requestIdMiddleware);
```

## 高级用法

### 1. 性能监控

```typescript
import { PerformanceLogger } from '@pkg/logger/examples';

const perfLogger = new PerformanceLogger('DatabaseService');

perfLogger.startTimer('userQuery');
// 执行数据库查询
perfLogger.endTimer('userQuery', { recordCount: 100 });
```

### 2. 结构化事件日志

```typescript
const eventLogger = createContextLogger('EventTracker');

eventLogger.info('User action', {
  event: 'button_click',
  userId: '123',
  buttonId: 'signup',
  metadata: {
    page: '/landing',
    campaign: 'summer_sale'
  }
});
```

### 3. 错误追踪

```typescript
const apiLogger = createContextLogger('ApiService');

try {
  await riskyOperation();
} catch (error) {
  apiLogger.error('API call failed', {
    operation: 'fetchUserData',
    userId: '123',
    error: {
      name: error.name,
      message: error.message,
      stack: error.stack
    }
  });
  
  throw error; // 重新抛出以便上层处理
}
```

## 环境配置

### 开发环境

- 控制台输出：彩色格式，易于阅读
- 日志级别：debug
- 文件输出：所有级别

### 生产环境

- 控制台输出：关闭
- 日志级别：info
- 文件输出：JSON 格式，便于分析
- 自动轮转：按日期和大小

### 测试环境

- HTTP 日志：跳过
- 最小化输出：避免干扰测试

## 最佳实践

### 1. 使用上下文日志

```typescript
// ✅ 推荐
const userLogger = createContextLogger('UserService');
userLogger.info('User created');

// ❌ 不推荐
logger.info('[UserService] User created');
```

### 2. 结构化数据

```typescript
// ✅ 推荐
logger.info('User login', {
  userId: '123',
  ip: '192.168.1.1',
  success: true
});

// ❌ 不推荐
logger.info(`User 123 logged in from 192.168.1.1 successfully`);
```

### 3. 错误处理

```typescript
// ✅ 推荐
try {
  await operation();
} catch (error) {
  logger.error('Operation failed', {
    operation: 'updateUser',
    userId: '123',
    error
  });
}

// ❌ 不推荐
catch (error) {
  logger.error(error.toString());
}
```

### 4. 性能敏感操作

```typescript
// ✅ 推荐 - 只在需要时记录详细信息
if (logger.isDebugEnabled()) {
  logger.debug('Detailed operation info', expensiveToCompute());
}

// 或使用惰性求值
logger.debug('Operation result', () => ({
  data: JSON.stringify(largeObject)
}));
```

### 5. 敏感信息

```typescript
// ✅ 推荐 - 过滤敏感信息
const safeUserData = {
  id: user.id,
  email: user.email,
  // 不记录密码等敏感信息
};
logger.info('User data processed', safeUserData);
```

## 监控和告警

### 日志分析

可以使用以下工具分析日志：

- **ELK Stack**: Elasticsearch + Logstash + Kibana
- **Fluentd**: 日志收集和转发
- **Grafana**: 日志可视化
- **DataDog**: 云端日志分析

### 告警设置

建议为以下情况设置告警：

- 错误日志频率异常
- 响应时间过长
- 特定错误模式
- 系统资源异常

## 故障排查

### 常见问题

1. **日志文件没有生成**
   - 检查 `logs/` 目录权限
   - 确认日志级别配置

2. **日志内容不完整**
   - 检查进程是否正常退出
   - 确认没有未捕获的异常

3. **性能影响**
   - 调整日志级别
   - 使用异步写入
   - 考虑日志采样

## 许可证

MIT
