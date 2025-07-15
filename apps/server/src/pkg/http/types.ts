import { AxiosRequestConfig, AxiosResponse } from "axios";

// HTTP 方法类型
export type HttpMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";

// 请求配置接口
export interface RequestConfig
  extends Omit<AxiosRequestConfig, "url" | "method"> {
  baseURL?: string;
  timeout?: number;
  retries?: number;
  retryDelay?: number;
}

// 响应接口
export interface ApiResponse<T = any> {
  data: T;
  status: number;
  statusText: string;
  headers: any;
}

// 错误响应接口
export interface ApiError {
  message: string;
  status?: number;
  code?: string;
  details?: any;
}

// 拦截器配置
export interface InterceptorConfig {
  requestInterceptor?: (
    config: AxiosRequestConfig
  ) => AxiosRequestConfig | Promise<AxiosRequestConfig>;
  responseInterceptor?: (
    response: AxiosResponse
  ) => AxiosResponse | Promise<AxiosResponse>;
  errorInterceptor?: (error: any) => Promise<any>;
}

// HTTP 客户端配置
export interface HttpClientConfig extends RequestConfig {
  interceptors?: InterceptorConfig;
}
