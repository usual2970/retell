import winston from "winston";
import DailyRotateFile from "winston-daily-rotate-file";
import path from "path";

// 日志级别定义
export enum LogLevel {
  ERROR = "error",
  WARN = "warn",
  INFO = "info",
  HTTP = "http",
  VERBOSE = "verbose",
  DEBUG = "debug",
  SILLY = "silly",
}

// 日志格式配置
const logFormat = winston.format.combine(
  winston.format.timestamp({
    format: "YYYY-MM-DD HH:mm:ss.SSS",
  }),
  winston.format.errors({ stack: true }),
  winston.format.json(),
  winston.format.printf(({ timestamp, level, message, stack, ...meta }) => {
    let log = `${timestamp} [${level.toUpperCase()}]: ${message}`;

    // 添加错误堆栈
    if (stack) {
      log += `\n${stack}`;
    }

    // 添加额外的元数据
    if (Object.keys(meta).length > 0) {
      log += `\n${JSON.stringify(meta, null, 2)}`;
    }

    return log;
  })
);

// 控制台格式（开发环境用）
const consoleFormat = winston.format.combine(
  winston.format.colorize(),
  winston.format.timestamp({
    format: "HH:mm:ss",
  }),
  winston.format.printf(({ timestamp, level, message, stack }) => {
    let log = `${timestamp} ${level}: ${message}`;
    if (stack) {
      log += `\n${stack}`;
    }
    return log;
  })
);

// 日志目录
const logDir = path.join(process.cwd(), "logs");

// 创建 Winston Logger
const createLogger = () => {
  const isDevelopment = process.env.NODE_ENV === "development";
  const logLevel = process.env.LOG_LEVEL || (isDevelopment ? "debug" : "info");

  const transports: winston.transport[] = [];

  // 控制台输出（开发环境）
  if (isDevelopment) {
    transports.push(
      new winston.transports.Console({
        level: logLevel,
        format: consoleFormat,
      })
    );
  }

  // 文件输出
  // 所有日志
  transports.push(
    new DailyRotateFile({
      filename: path.join(logDir, "application-%DATE%.log"),
      datePattern: "YYYY-MM-DD",
      maxSize: "20m",
      maxFiles: "14d",
      level: logLevel,
      format: logFormat,
    })
  );

  // 错误日志单独存储
  transports.push(
    new DailyRotateFile({
      filename: path.join(logDir, "error-%DATE%.log"),
      datePattern: "YYYY-MM-DD",
      maxSize: "20m",
      maxFiles: "30d",
      level: "error",
      format: logFormat,
    })
  );

  return winston.createLogger({
    level: logLevel,
    format: logFormat,
    transports,
    // 处理未捕获的异常
    exceptionHandlers: [
      new DailyRotateFile({
        filename: path.join(logDir, "exceptions-%DATE%.log"),
        datePattern: "YYYY-MM-DD",
        maxSize: "20m",
        maxFiles: "30d",
      }),
    ],
    // 处理未捕获的 Promise 拒绝
    rejectionHandlers: [
      new DailyRotateFile({
        filename: path.join(logDir, "rejections-%DATE%.log"),
        datePattern: "YYYY-MM-DD",
        maxSize: "20m",
        maxFiles: "30d",
      }),
    ],
  });
};

// 导出 logger 实例
export const logger = createLogger();

// 创建带上下文的 logger
export class ContextLogger {
  private context: string;
  private logger: winston.Logger;

  constructor(context: string) {
    this.context = context;
    this.logger = logger;
  }

  private formatMessage(message: string): string {
    return `[${this.context}] ${message}`;
  }

  error(message: string, error?: Error | object): void {
    this.logger.error(this.formatMessage(message), error);
  }

  warn(message: string, meta?: object): void {
    this.logger.warn(this.formatMessage(message), meta);
  }

  info(message: string, meta?: object): void {
    this.logger.info(this.formatMessage(message), meta);
  }

  http(message: string, meta?: object): void {
    this.logger.http(this.formatMessage(message), meta);
  }

  debug(message: string, meta?: object): void {
    this.logger.debug(this.formatMessage(message), meta);
  }

  verbose(message: string, meta?: object): void {
    this.logger.verbose(this.formatMessage(message), meta);
  }
}

// 工具函数：创建带上下文的 logger
export const createContextLogger = (context: string): ContextLogger => {
  return new ContextLogger(context);
};
