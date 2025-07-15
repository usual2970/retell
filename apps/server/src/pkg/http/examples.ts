import { HttpClient, HttpUtils, http, HTTP_STATUS } from "./index";

/**
 * HTTP 库使用示例
 */

// 示例 1: 使用默认客户端进行简单请求
export async function example1() {
  try {
    // GET 请求
    const response = await http.get<{ message: string }>(
      "https://api.example.com/users"
    );
    console.log("Users:", response.data);

    // POST 请求
    const createResponse = await http.post("https://api.example.com/users", {
      name: "John Doe",
      email: "john@example.com",
    });
    console.log("Created user:", createResponse.data);
  } catch (error) {
    console.error("Request failed:", error);
  }
}

// 示例 2: 创建自定义客户端
export async function example2() {
  const apiClient = new HttpClient({
    baseURL: "https://api.example.com",
    timeout: 15000,
    retries: 5,
    headers: {
      "X-API-Key": "your-api-key",
    },
  });

  try {
    const response = await apiClient.get("/users/1");
    console.log("User details:", response.data);
  } catch (error) {
    console.error("API request failed:", error);
  }
}

// 示例 3: 使用认证客户端
export async function example3() {
  const authClient = HttpUtils.createAuthenticatedClient("your-jwt-token", {
    baseURL: "https://api.example.com",
    timeout: 20000,
  });

  try {
    const profile = await authClient.get("/profile");
    console.log("User profile:", profile.data);

    const updateResponse = await authClient.put("/profile", {
      name: "Updated Name",
    });
    console.log("Profile updated:", updateResponse.data);
  } catch (error) {
    console.error("Authenticated request failed:", error);
  }
}

// 示例 4: 使用拦截器
export async function example4() {
  const clientWithInterceptors = new HttpClient({
    baseURL: "https://api.example.com",
    interceptors: {
      requestInterceptor: (config) => {
        console.log("Sending request to:", config.url);
        // 添加时间戳
        config.params = { ...config.params, _t: Date.now() };
        return config;
      },
      responseInterceptor: (response) => {
        console.log("Response received:", response.status);
        return response;
      },
      errorInterceptor: async (error) => {
        console.error("Request error:", error.message);

        // 如果是 401 错误，可能需要刷新 token
        if (error.status === HTTP_STATUS.UNAUTHORIZED) {
          console.log("Attempting to refresh token...");
          // 在这里处理 token 刷新逻辑
        }

        throw error;
      },
    },
  });

  try {
    const response = await clientWithInterceptors.get("/protected-resource");
    console.log("Protected data:", response.data);
  } catch (error) {
    console.error("Protected request failed:", error);
  }
}

// 示例 5: 文件上传
export async function example5() {
  const uploadClient = new HttpClient({
    baseURL: "https://api.example.com",
    timeout: 60000, // 文件上传需要更长的超时时间
  });

  const formData = new FormData();
  // formData.append('file', fileInput); // 在浏览器环境中使用
  formData.append("description", "My uploaded file");

  try {
    const response = await uploadClient.post("/upload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
      // 可以监听上传进度
      onUploadProgress: (progressEvent) => {
        const percentCompleted = Math.round(
          (progressEvent.loaded * 100) / (progressEvent.total || 1)
        );
        console.log(`Upload progress: ${percentCompleted}%`);
      },
    });
    console.log("File uploaded:", response.data);
  } catch (error) {
    console.error("Upload failed:", error);
  }
}

// 示例 6: 批量请求和错误处理
export async function example6() {
  const apiClient = HttpUtils.createApiClient("https://api.example.com");

  const userIds = [1, 2, 3, 4, 5];

  try {
    // 并发请求
    const userPromises = userIds.map((id) =>
      apiClient.get(`/users/${id}`).catch((error) => ({ error, id }))
    );

    const results = await Promise.all(userPromises);

    // 处理结果
    const successResults = results.filter((result) => !("error" in result));
    const errorResults = results.filter((result) => "error" in result);

    console.log(`Successfully fetched ${successResults.length} users`);
    console.log(`Failed to fetch ${errorResults.length} users`);

    errorResults.forEach((result) => {
      if ("error" in result) {
        console.error(
          `Failed to fetch user ${result.id}:`,
          result.error.message
        );
      }
    });
  } catch (error) {
    console.error("Batch request failed:", error);
  }
}

// 示例 7: 使用工具函数
export function example7() {
  // 构建带参数的 URL
  const url = HttpUtils.buildUrl("https://api.example.com/search", {
    q: "javascript",
    page: 1,
    limit: 10,
  });
  console.log("Search URL:", url);

  // 检查状态码
  const status = 404;
  if (HttpUtils.isClientError(status)) {
    console.log("This is a client error");
  }

  // 格式化文件大小
  const fileSize = HttpUtils.formatFileSize(1024 * 1024 * 2.5);
  console.log("File size:", fileSize); // "2.5 MB"
}
