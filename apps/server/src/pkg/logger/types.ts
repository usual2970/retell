import { Request } from "express";

// 扩展 Express Request 类型
declare global {
  namespace Express {
    interface Request {
      id?: string;
      startTime?: number;
      user?: {
        id: string;
        [key: string]: any;
      };
    }
  }
}
