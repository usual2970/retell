import { HttpClient } from "./client";
import { HttpClientConfig } from "./types";

/**
 * HTTP 工具类 - 提供常用的 HTTP 操作方法
 */
export class HttpUtils {
  /**
   * 创建一个新的 HTTP 客户端实例
   */
  static createClient(config?: HttpClientConfig): HttpClient {
    return new HttpClient(config);
  }

  /**
   * 创建一个带有认证的 HTTP 客户端
   */
  static createAuthenticatedClient(
    token: string,
    config?: HttpClientConfig
  ): HttpClient {
    const client = new HttpClient(config);
    client.setAuthToken(token);
    return client;
  }

  /**
   * 创建一个用于 API 调用的客户端
   */
  static createApiClient(
    baseURL: string,
    config?: Omit<HttpClientConfig, "baseURL">
  ): HttpClient {
    return new HttpClient({
      baseURL,
      timeout: 30000,
      ...config,
    });
  }

  /**
   * URL 编码参数
   */
  static encodeParams(params: Record<string, any>): string {
    return Object.keys(params)
      .filter((key) => params[key] !== undefined && params[key] !== null)
      .map(
        (key) => `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`
      )
      .join("&");
  }

  /**
   * 构建带查询参数的 URL
   */
  static buildUrl(baseUrl: string, params?: Record<string, any>): string {
    if (!params || Object.keys(params).length === 0) {
      return baseUrl;
    }

    const queryString = this.encodeParams(params);
    const separator = baseUrl.includes("?") ? "&" : "?";
    return `${baseUrl}${separator}${queryString}`;
  }

  /**
   * 检查响应是否成功
   */
  static isSuccessResponse(status: number): boolean {
    return status >= 200 && status < 300;
  }

  /**
   * 检查是否为客户端错误
   */
  static isClientError(status: number): boolean {
    return status >= 400 && status < 500;
  }

  /**
   * 检查是否为服务器错误
   */
  static isServerError(status: number): boolean {
    return status >= 500 && status < 600;
  }

  /**
   * 格式化文件大小
   */
  static formatFileSize(bytes: number): string {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  }
}

/**
 * 默认的 HTTP 客户端实例
 */
export const defaultHttpClient = new HttpClient({
  timeout: 10000,
  retries: 3,
  retryDelay: 1000,
});

/**
 * 快捷方法 - 直接使用默认客户端
 */
export const http = {
  get: <T = any>(url: string, config?: any) =>
    defaultHttpClient.get<T>(url, config),
  post: <T = any>(url: string, data?: any, config?: any) =>
    defaultHttpClient.post<T>(url, data, config),
  put: <T = any>(url: string, data?: any, config?: any) =>
    defaultHttpClient.put<T>(url, data, config),
  delete: <T = any>(url: string, config?: any) =>
    defaultHttpClient.delete<T>(url, config),
  patch: <T = any>(url: string, data?: any, config?: any) =>
    defaultHttpClient.patch<T>(url, data, config),
};
