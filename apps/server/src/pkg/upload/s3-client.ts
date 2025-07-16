import { S3Client, PutObjectCommand } from "@aws-sdk/client-s3";
import type { S3Config, UploadResult, FileInput } from "./types";

export class S3UploadClient {
  private client: S3Client;
  private bucket: string;
  private domain: string;

  constructor(config: S3Config) {
    this.bucket = config.bucket;
    this.domain = config.domain;
    this.client = new S3Client({
      region: config.region,
      endpoint: config.endpoint,
      credentials: {
        accessKeyId: config.accessKeyId,
        secretAccessKey: config.secretAccessKey,
      },
    });

    console.log("S3UploadClient initialized:", {
      region: config.region,
      bucket: config.bucket,
    });
  }

  /**
   * 上传单个文件到 S3
   */
  async upload(
    file: FileInput,
    key: string,
    options?: {
      contentType?: string;
      metadata?: Record<string, string>;
    }
  ): Promise<UploadResult> {
    const requestId = `upload-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

    console.log("Starting file upload:", {
      requestId,
      key,
      bucket: this.bucket,
      contentType: options?.contentType,
    });

    try {
      // 转换文件为 Buffer
      const buffer = typeof file === "string" ? Buffer.from(file) : file;

      // 创建上传命令
      const command = new PutObjectCommand({
        Bucket: this.bucket,
        Key: key,
        Body: buffer,
        ContentType: options?.contentType || "application/octet-stream",
        Metadata: options?.metadata,
      });

      // 执行上传
      const response = await this.client.send(command);

      const result: UploadResult = {
        key,
        url: `${this.domain}/${key}`,
        size: buffer.length,
      };

      console.log("File upload completed:", {
        requestId,
        key,
        url: result.url,
        size: result.size,
        etag: response.ETag,
      });

      return result;
    } catch (error) {
      console.error("File upload failed:", {
        requestId,
        key,
        bucket: this.bucket,
        error: error instanceof Error ? error.message : String(error),
      });
      throw error;
    }
  }

  /**
   * 生成文件的公共访问 URL
   */
  getPublicUrl(key: string): string {
    return `https://${this.bucket}.s3.amazonaws.com/${key}`;
  }
}
