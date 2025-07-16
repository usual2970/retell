import { S3UploadClient } from "./s3-client";
import { FileInput, UploadResult } from "./types";

// S3 上传库主入口
export { S3UploadClient } from "./s3-client";
export type { S3Config, UploadResult, FileInput } from "./types";

const s3UploadClient = new S3UploadClient({
  region: process.env.AWS_REGION || "us-east-1",
  accessKeyId: process.env.AWS_ACCESS_KEY_ID || "",
  secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY || "",
  bucket: process.env.AWS_S3_BUCKET || "",
  domain: process.env.AWS_S3_DOMAIN || "",
  endpoint: process.env.AWS_S3_ENDPOINT || "",
});

export const upload = async (
  file: FileInput,
  key: string,
  options?: {
    contentType?: string;
    metadata?: Record<string, string>;
  }
): Promise<UploadResult> => {
  return await s3UploadClient.upload(file, key, options);
};
