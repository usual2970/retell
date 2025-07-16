import morgan from "morgan";
import { Request, Response, RequestHandler } from "express";
import { logger } from "./index";

// 创建自定义 token
morgan.token("id", (req: Request) => {
  return req.id || "unknown";
});

morgan.token("user-id", (req: Request) => {
  return (req as any).user?.id || "anonymous";
});

morgan.token("body", (req: Request) => {
  // 只记录非敏感数据
  if (req.method === "POST" || req.method === "PUT" || req.method === "PATCH") {
    const body = { ...req.body };
    // 移除敏感字段
    delete body.password;
    delete body.token;
    delete body.secret;
    return JSON.stringify(body);
  }
  return "";
});

morgan.token("response-time-ms", (req: Request, res: Response) => {
  const start = (req as any).startTime;
  if (start) {
    return `${Date.now() - start}ms`;
  }
  return "";
});

// 开发环境格式
const developmentFormat =
  ":method :url :status :res[content-length] - :response-time";

// 生产环境格式
const productionFormat = JSON.stringify({
  timestamp: ":date[iso]",
  method: ":method",
  url: ":url",
  status: ":status",
  contentLength: ":res[content-length]",
  responseTime: ":response-time",
  remoteAddr: ":remote-addr",
  userAgent: ":user-agent",
  userId: ":user-id",
  requestId: ":id",
  body: ":body",
});

// HTTP 日志流
const httpLogStream = {
  write: (message: string) => {
    // 移除末尾的换行符
    const cleanMessage = message.trim();

    try {
      // 尝试解析为 JSON（生产环境）
      const logData = JSON.parse(cleanMessage);
      logger.http("HTTP Request", logData);
    } catch {
      // 普通字符串格式（开发环境）
      logger.http(cleanMessage);
    }
  },
};

// 跳过某些请求的日志记录
const skip = (req: Request, res: Response) => {
  // 跳过健康检查和静态资源
  if (req.url === "/health" || req.url === "/favicon.ico") {
    return true;
  }

  // 在测试环境中跳过所有日志
  if (process.env.NODE_ENV === "test") {
    return true;
  }

  return false;
};

// 创建 Morgan 中间件
export const createMorganMiddleware = (): RequestHandler => {
  const isDevelopment = process.env.NODE_ENV === "development";

  return morgan(isDevelopment ? developmentFormat : productionFormat, {
    stream: httpLogStream,
    skip,
  });
};

// 请求 ID 中间件
export const requestIdMiddleware = (
  req: Request,
  res: Response,
  next: Function
) => {
  // 生成请求 ID
  req.id =
    (req.headers["x-request-id"] as string) ||
    `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

  // 记录请求开始时间
  (req as any).startTime = Date.now();

  // 设置响应头
  res.setHeader("X-Request-ID", req.id);

  next();
};
