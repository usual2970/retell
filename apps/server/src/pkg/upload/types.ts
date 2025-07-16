// S3 上传配置类型
export interface S3Config {
  region: string;
  accessKeyId: string;
  secretAccessKey: string;
  bucket: string;
  domain: string;
  endpoint: string;
}

// 上传结果
export interface UploadResult {
  key: string;
  url: string;
  size: number;
}

// 文件输入类型
export type FileInput = Buffer | string;
