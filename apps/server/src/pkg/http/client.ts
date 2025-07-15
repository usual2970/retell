import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios';
import { HttpClientConfig, ApiResponse, ApiError, RequestConfig } from './types';

export class HttpClient {
  private instance: AxiosInstance;
  private defaultRetries: number = 3;
  private defaultRetryDelay: number = 1000;

  constructor(config: HttpClientConfig = {}) {
    // 创建 axios 实例
    this.instance = axios.create({
      baseURL: config.baseURL || '',
      timeout: config.timeout || 10000,
      headers: {
        'Content-Type': 'application/json',
        ...config.headers,
      },
      ...config,
    });

    this.defaultRetries = config.retries || this.defaultRetries;
    this.defaultRetryDelay = config.retryDelay || this.defaultRetryDelay;

    // 设置拦截器
    this.setupInterceptors(config);
  }

  /**
   * 设置请求和响应拦截器
   */
  private setupInterceptors(config: HttpClientConfig) {
    // 请求拦截器
    this.instance.interceptors.request.use(
      (requestConfig: InternalAxiosRequestConfig) => {
        // 添加请求时间戳
        requestConfig.metadata = { startTime: Date.now() };
        
        // 自定义请求拦截器
        if (config.interceptors?.requestInterceptor) {
          const result = config.interceptors.requestInterceptor(requestConfig);
          return result as InternalAxiosRequestConfig;
        }
        
        return requestConfig;
      },
      (error) => {
        return Promise.reject(this.formatError(error));
      }
    );

    // 响应拦截器
    this.instance.interceptors.response.use(
      (response) => {
        // 记录响应时间
        const startTime = response.config.metadata?.startTime;
        if (startTime) {
          const duration = Date.now() - startTime;
          console.log(`API ${response.config.method?.toUpperCase()} ${response.config.url} - ${duration}ms`);
        }

        // 自定义响应拦截器
        if (config.interceptors?.responseInterceptor) {
          return config.interceptors.responseInterceptor(response);
        }

        return response;
      },
      async (error) => {
        // 自定义错误拦截器
        if (config.interceptors?.errorInterceptor) {
          return config.interceptors.errorInterceptor(error);
        }

        return Promise.reject(this.formatError(error));
      }
    );
  }

  /**
   * 格式化错误信息
   */
  private formatError(error: any): ApiError {
    if (error.response) {
      // 服务器响应了错误状态码
      return {
        message: error.response.data?.message || error.message || 'Request failed',
        status: error.response.status,
        code: error.response.data?.code || error.code,
        details: error.response.data,
      };
    } else if (error.request) {
      // 请求已发出但没有收到响应
      return {
        message: 'Network error or server not responding',
        code: 'NETWORK_ERROR',
        details: error.request,
      };
    } else {
      // 其他错误
      return {
        message: error.message || 'Unknown error occurred',
        code: 'UNKNOWN_ERROR',
        details: error,
      };
    }
  }

  /**
   * 重试机制
   */
  private async executeWithRetry<T>(
    requestFn: () => Promise<AxiosResponse<T>>,
    retries: number = this.defaultRetries
  ): Promise<AxiosResponse<T>> {
    try {
      return await requestFn();
    } catch (error: any) {
      if (retries > 0 && this.shouldRetry(error)) {
        console.log(`Request failed, retrying... (${retries} attempts left)`);
        await this.delay(this.defaultRetryDelay);
        return this.executeWithRetry(requestFn, retries - 1);
      }
      throw error;
    }
  }

  /**
   * 判断是否应该重试
   */
  private shouldRetry(error: any): boolean {
    // 网络错误或5xx服务器错误时重试
    return (
      !error.response ||
      (error.response.status >= 500 && error.response.status < 600) ||
      error.code === 'NETWORK_ERROR' ||
      error.code === 'ECONNABORTED'
    );
  }

  /**
   * 延迟函数
   */
  private delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  /**
   * GET 请求
   */
  async get<T = any>(url: string, config?: RequestConfig): Promise<ApiResponse<T>> {
    const response = await this.executeWithRetry<T>(
      () => this.instance.get<T>(url, config),
      config?.retries
    );
    
    return {
      data: response.data,
      status: response.status,
      statusText: response.statusText,
      headers: response.headers,
    };
  }

  /**
   * POST 请求
   */
  async post<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>> {
    const response = await this.executeWithRetry<T>(
      () => this.instance.post<T>(url, data, config),
      config?.retries
    );
    
    return {
      data: response.data,
      status: response.status,
      statusText: response.statusText,
      headers: response.headers,
    };
  }

  /**
   * PUT 请求
   */
  async put<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>> {
    const response = await this.executeWithRetry<T>(
      () => this.instance.put<T>(url, data, config),
      config?.retries
    );
    
    return {
      data: response.data,
      status: response.status,
      statusText: response.statusText,
      headers: response.headers,
    };
  }

  /**
   * DELETE 请求
   */
  async delete<T = any>(url: string, config?: RequestConfig): Promise<ApiResponse<T>> {
    const response = await this.executeWithRetry<T>(
      () => this.instance.delete<T>(url, config),
      config?.retries
    );
    
    return {
      data: response.data,
      status: response.status,
      statusText: response.statusText,
      headers: response.headers,
    };
  }

  /**
   * PATCH 请求
   */
  async patch<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>> {
    const response = await this.executeWithRetry<T>(
      () => this.instance.patch<T>(url, data, config),
      config?.retries
    );
    
    return {
      data: response.data,
      status: response.status,
      statusText: response.statusText,
      headers: response.headers,
    };
  }

  /**
   * 设置默认请求头
   */
  setDefaultHeader(key: string, value: string): void {
    this.instance.defaults.headers.common[key] = value;
  }

  /**
   * 移除默认请求头
   */
  removeDefaultHeader(key: string): void {
    delete this.instance.defaults.headers.common[key];
  }

  /**
   * 设置认证 token
   */
  setAuthToken(token: string, type: string = 'Bearer'): void {
    this.setDefaultHeader('Authorization', `${type} ${token}`);
  }

  /**
   * 清除认证 token
   */
  clearAuthToken(): void {
    this.removeDefaultHeader('Authorization');
  }

  /**
   * 获取 axios 实例（用于高级用法）
   */
  getAxiosInstance(): AxiosInstance {
    return this.instance;
  }
}

// 扩展 InternalAxiosRequestConfig 类型以支持 metadata
declare module 'axios' {
  interface InternalAxiosRequestConfig {
    metadata?: {
      startTime?: number;
    };
  }
}
