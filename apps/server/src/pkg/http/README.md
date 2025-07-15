# HTTP 库

一个基于 axios 封装的 TypeScript HTTP 客户端库，提供了丰富的功能和易用的 API。

## 特性

- 🚀 基于 axios，性能优异
- 📘 完整的 TypeScript 支持
- 🔄 自动重试机制
- ⏱️ 请求超时控制
- 🎯 拦截器支持
- 📊 请求响应时间统计
- 🛡️ 统一错误处理
- 🔐 认证支持
- 🛠️ 丰富的工具函数

## 安装

```bash
pnpm add axios
```

## 基本用法

### 1. 使用默认客户端

```typescript
import { http } from '@pkg/http';

// GET 请求
const response = await http.get<User[]>('/api/users');
console.log(response.data);

// POST 请求
const newUser = await http.post<User>('/api/users', {
  name: 'John Doe',
  email: 'john@example.com'
});
```

### 2. 创建自定义客户端

```typescript
import { HttpClient } from '@pkg/http';

const apiClient = new HttpClient({
  baseURL: 'https://api.example.com',
  timeout: 15000,
  retries: 3,
  headers: {
    'X-API-Key': 'your-api-key'
  }
});

const data = await apiClient.get('/users');
```

### 3. 认证客户端

```typescript
import { HttpUtils } from '@pkg/http';

const authClient = HttpUtils.createAuthenticatedClient('your-jwt-token', {
  baseURL: 'https://api.example.com'
});

const profile = await authClient.get('/profile');
```

## 高级功能

### 拦截器

```typescript
const client = new HttpClient({
  baseURL: 'https://api.example.com',
  interceptors: {
    requestInterceptor: (config) => {
      // 添加时间戳
      config.params = { ...config.params, _t: Date.now() };
      return config;
    },
    responseInterceptor: (response) => {
      console.log('Response time:', response.headers['x-response-time']);
      return response;
    },
    errorInterceptor: async (error) => {
      if (error.status === 401) {
        // 处理认证失效
        await refreshToken();
      }
      throw error;
    }
  }
});
```

### 重试机制

```typescript
const response = await apiClient.get('/unstable-endpoint', {
  retries: 5,          // 重试 5 次
  retryDelay: 2000     // 重试间隔 2 秒
});
```

### 文件上传

```typescript
const formData = new FormData();
formData.append('file', file);

const response = await apiClient.post('/upload', formData, {
  headers: { 'Content-Type': 'multipart/form-data' },
  timeout: 60000,
  onUploadProgress: (progressEvent) => {
    const progress = Math.round(
      (progressEvent.loaded * 100) / progressEvent.total
    );
    console.log(\`Upload progress: \${progress}%\`);
  }
});
```

## API 文档

### HttpClient

#### 构造函数
```typescript
new HttpClient(config?: HttpClientConfig)
```

#### 方法
- `get<T>(url: string, config?: RequestConfig): Promise<ApiResponse<T>>`
- `post<T>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>>`
- `put<T>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>>`
- `delete<T>(url: string, config?: RequestConfig): Promise<ApiResponse<T>>`
- `patch<T>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>>`
- `setAuthToken(token: string, type?: string): void`
- `clearAuthToken(): void`
- `setDefaultHeader(key: string, value: string): void`
- `removeDefaultHeader(key: string): void`

### HttpUtils

静态工具方法：

- `createClient(config?: HttpClientConfig): HttpClient`
- `createAuthenticatedClient(token: string, config?: HttpClientConfig): HttpClient`
- `createApiClient(baseURL: string, config?: HttpClientConfig): HttpClient`
- `buildUrl(baseUrl: string, params?: Record<string, any>): string`
- `encodeParams(params: Record<string, any>): string`
- `isSuccessResponse(status: number): boolean`
- `isClientError(status: number): boolean`
- `isServerError(status: number): boolean`
- `formatFileSize(bytes: number): string`

### 类型定义

```typescript
interface ApiResponse<T> {
  data: T;
  status: number;
  statusText: string;
  headers: any;
}

interface ApiError {
  message: string;
  status?: number;
  code?: string;
  details?: any;
}

interface HttpClientConfig {
  baseURL?: string;
  timeout?: number;
  retries?: number;
  retryDelay?: number;
  headers?: Record<string, string>;
  interceptors?: InterceptorConfig;
}
```

## 常量

### HTTP 状态码
```typescript
import { HTTP_STATUS } from '@pkg/http';

if (error.status === HTTP_STATUS.UNAUTHORIZED) {
  // 处理未授权错误
}
```

### HTTP 方法
```typescript
import { HTTP_METHODS } from '@pkg/http';

console.log(HTTP_METHODS.GET); // "GET"
```

## 错误处理

库提供了统一的错误格式：

```typescript
try {
  const response = await apiClient.get('/api/data');
} catch (error: ApiError) {
  console.log('Error message:', error.message);
  console.log('Status code:', error.status);
  console.log('Error code:', error.code);
  console.log('Error details:', error.details);
}
```

## 最佳实践

1. **使用 TypeScript 类型**：为请求和响应数据定义明确的类型
2. **配置合理的超时时间**：根据 API 响应时间设置合适的超时
3. **使用重试机制**：对于可能不稳定的网络环境
4. **添加拦截器**：统一处理认证、日志、错误等
5. **创建专用客户端**：为不同的 API 服务创建专用的客户端实例

## 许可证

MIT
