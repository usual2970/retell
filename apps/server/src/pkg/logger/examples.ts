import { logger, createContextLogger, ContextLogger } from "./index";
import { createMorganMiddleware, requestIdMiddleware } from "./morgan";

/**
 * 日志使用示例
 */

// 示例 1: 基础日志使用
export function basicLoggingExample() {
  // 不同级别的日志
  logger.error("This is an error message");
  logger.warn("This is a warning message");
  logger.info("This is an info message");
  logger.http("This is an HTTP message");
  logger.debug("This is a debug message");

  // 带有额外数据的日志
  logger.info("User login", {
    userId: "12345",
    email: "user@example.com",
    ip: "192.168.1.1",
  });

  // 错误日志
  try {
    throw new Error("Something went wrong");
  } catch (error) {
    logger.error("Operation failed", error);
  }
}

// 示例 2: 上下文日志使用
export function contextLoggingExample() {
  const userLogger = createContextLogger("UserService");
  const authLogger = createContextLogger("AuthService");

  userLogger.info("Creating new user", { email: "test@example.com" });
  authLogger.warn("Invalid login attempt", { ip: "192.168.1.100" });
}

// 示例 3: 在 Express 中使用
export function expressLoggingExample() {
  const express = require("express");
  const app = express();

  // 添加请求 ID 中间件
  app.use(requestIdMiddleware);

  // 添加 Morgan HTTP 日志中间件
  app.use(createMorganMiddleware());

  // 在路由中使用日志
  app.get("/users/:id", (req: any, res: any) => {
    const userLogger = createContextLogger("UserController");

    userLogger.info("Fetching user", {
      requestId: req.id,
      userId: req.params.id,
    });

    try {
      // 模拟业务逻辑
      const user = { id: req.params.id, name: "John Doe" };

      userLogger.info("User fetched successfully", {
        requestId: req.id,
        userId: user.id,
      });

      res.json(user);
    } catch (error) {
      userLogger.error("Failed to fetch user", {
        requestId: req.id,
        userId: req.params.id,
        error,
      });

      res.status(500).json({ error: "Internal server error" });
    }
  });

  return app;
}

// 示例 4: 异步操作日志
export async function asyncLoggingExample() {
  const dbLogger = createContextLogger("Database");

  try {
    dbLogger.info("Starting database connection");

    // 模拟异步数据库操作
    await new Promise((resolve) => setTimeout(resolve, 1000));

    dbLogger.info("Database connected successfully");

    // 模拟查询
    dbLogger.debug("Executing query", {
      sql: "SELECT * FROM users WHERE active = ?",
      params: [true],
    });

    await new Promise((resolve) => setTimeout(resolve, 500));

    dbLogger.info("Query executed successfully", {
      rowCount: 42,
      executionTime: "523ms",
    });
  } catch (error) {
    dbLogger.error("Database operation failed", error as Error);
  }
}

// 示例 5: 性能监控日志
export class PerformanceLogger {
  private logger: ContextLogger;
  private timers: Map<string, number> = new Map();

  constructor(context: string) {
    this.logger = createContextLogger(context);
  }

  startTimer(operation: string): void {
    this.timers.set(operation, Date.now());
    this.logger.debug(`Started ${operation}`);
  }

  endTimer(operation: string, metadata?: object): void {
    const startTime = this.timers.get(operation);
    if (startTime) {
      const duration = Date.now() - startTime;
      this.timers.delete(operation);

      this.logger.info(`Completed ${operation}`, {
        duration: `${duration}ms`,
        ...metadata,
      });
    } else {
      this.logger.warn(`Timer not found for operation: ${operation}`);
    }
  }
}

// 使用性能日志的示例
export async function performanceExample() {
  const perfLogger = new PerformanceLogger("ApiService");

  perfLogger.startTimer("fetchUserData");

  try {
    // 模拟 API 调用
    await new Promise((resolve) => setTimeout(resolve, 2000));

    perfLogger.endTimer("fetchUserData", {
      userCount: 100,
      cacheHit: true,
    });
  } catch (error) {
    perfLogger.endTimer("fetchUserData");
    throw error;
  }
}

// 示例 6: 结构化日志记录
export function structuredLoggingExample() {
  const apiLogger = createContextLogger("ApiHandler");

  // 记录 API 请求
  apiLogger.info("API request received", {
    endpoint: "/api/users",
    method: "POST",
    contentType: "application/json",
    userAgent: "MyApp/1.0",
    timestamp: new Date().toISOString(),
  });

  // 记录业务事件
  apiLogger.info("User created", {
    event: "user_created",
    userId: "user_123",
    email: "user@example.com",
    source: "registration_form",
    metadata: {
      referrer: "google",
      campaign: "signup_bonus",
    },
  });

  // 记录系统指标
  apiLogger.info("System metrics", {
    type: "metrics",
    cpu_usage: 45.2,
    memory_usage: 67.8,
    active_connections: 150,
    response_time_avg: 234,
  });
}
