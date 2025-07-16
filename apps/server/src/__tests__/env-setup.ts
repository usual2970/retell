import dotenv from 'dotenv';
import path from 'path';

// 加载测试环境变量
dotenv.config({ path: path.resolve(__dirname, '../../.env.test') });

// 如果没有测试环境文件，则加载开发环境文件
if (!process.env.DATABASE_URL) {
  dotenv.config({ path: path.resolve(__dirname, '../../.env') });
}

// 设置测试环境特定的变量
process.env.NODE_ENV = 'test';