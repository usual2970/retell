package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/pkg/logger"
)

// RequestLogConfig 请求日志中间件配置
type RequestLogConfig struct {
	// 要排除的路径前缀，这些路径不会被记录
	SkipPaths []string

	// 是否记录请求体，对于大型请求体可能需要禁用
	LogRequestBody bool

	// 是否记录响应体，可能包含敏感信息需要禁用
	LogResponseBody bool

	// 最大记录的请求/响应体大小（字节）
	MaxBodySize int

	// 敏感字段列表，这些字段会在日志中被替换为 "***"
	SensitiveFields []string
}

// DefaultRequestLogConfig 返回默认请求日志配置
func DefaultRequestLogConfig() RequestLogConfig {
	return RequestLogConfig{
		SkipPaths:       []string{"/health", "/metrics"},
		LogRequestBody:  true,
		LogResponseBody: true,
		MaxBodySize:     10 * 1024, // 10KB
		SensitiveFields: []string{"password", "token", "secret", "authorization", "api_key"},
	}
}

// 自定义响应体写入器
type responseBodyWriter struct {
	io.Writer
	http.ResponseWriter
	statusCode int
}

// WriteHeader 实现 ResponseWriter 接口
func (r *responseBodyWriter) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Write 实现 ResponseWriter 接口
func (r *responseBodyWriter) Write(b []byte) (int, error) {
	return r.Writer.Write(b)
}

// Status 返回状态码
func (r *responseBodyWriter) Status() int {
	return r.statusCode
}

// RequestLogger 是请求日志中间件函数
func RequestLogger() echo.MiddlewareFunc {
	return RequestLoggerWithConfig(DefaultRequestLogConfig())
}

// RequestLoggerWithConfig 使用自定义配置的请求日志中间件
func RequestLoggerWithConfig(config RequestLogConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 检查是否跳过该路径
			path := c.Request().URL.Path
			for _, skipPath := range config.SkipPaths {
				if strings.HasPrefix(path, skipPath) {
					return next(c)
				}
			}

			start := time.Now()

			// 获取请求信息
			req := c.Request()
			method := req.Method
			userAgent := req.UserAgent()
			ip := c.RealIP()

			// 获取并记录请求体
			var reqBody string
			if config.LogRequestBody && req.Body != nil && req.Method != http.MethodGet {
				// 读取请求体
				var bodyBytes []byte
				if req.Body != nil {
					bodyBytes, _ = io.ReadAll(req.Body)
					// 限制请求体大小
					if len(bodyBytes) > config.MaxBodySize {
						reqBody = string(bodyBytes[:config.MaxBodySize]) + "...(truncated)"
					} else {
						reqBody = string(bodyBytes)
					}
					// 替换敏感字段
					reqBody = maskSensitiveData(reqBody, config.SensitiveFields)
					// 恢复请求体以便后续处理
					req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}

			// 捕获响应体
			var resBody []byte
			if config.LogResponseBody {
				// 创建响应体捕获器
				resBodyBuffer := bytes.NewBuffer(nil)
				writer := &responseBodyWriter{
					Writer:         io.MultiWriter(c.Response().Writer, resBodyBuffer),
					ResponseWriter: c.Response().Writer,
					statusCode:     http.StatusOK, // 默认状态码
				}
				c.Response().Writer = writer

				// 调用下一个处理函数
				err := next(c)

				// 获取响应体和状态码
				if resBodyBuffer.Len() > 0 {
					if resBodyBuffer.Len() > config.MaxBodySize {
						resBody = append(resBodyBuffer.Bytes()[:config.MaxBodySize], []byte("...(truncated)")...)
					} else {
						resBody = resBodyBuffer.Bytes()
					}
				}

				// 生成日志
				statusCode := writer.Status()
				latency := time.Since(start)

				// 构建日志字段
				fields := map[string]any{
					"status_code":   statusCode,
					"latency":       latency.String(),
					"latency_ms":    float64(latency.Nanoseconds()) / 1e6,
					"ip":            ip,
					"method":        method,
					"path":          path,
					"user_agent":    userAgent,
					"query_params":  req.URL.RawQuery,
					"request_id":    c.Response().Header().Get(echo.HeaderXRequestID),
					"request_body":  reqBody,
					"response_body": maskSensitiveData(string(resBody), config.SensitiveFields),
				}

				// 根据状态码决定日志级别
				log := logger.WithFields(fields)

				if statusCode >= 500 {
					log.Error("Server error")
				} else if statusCode >= 400 {
					log.Warn("Client error")
				} else {
					log.Info("Request completed")
				}

				return err
			} else {
				// 如果不记录响应体，则先调用下一个处理函数
				err := next(c)

				// 生成日志
				statusCode := c.Response().Status
				latency := time.Since(start)

				// 构建日志字段
				fields := map[string]interface{}{
					"status_code":  statusCode,
					"latency":      latency.String(),
					"latency_ms":   float64(latency.Nanoseconds()) / 1e6,
					"ip":           ip,
					"method":       method,
					"path":         path,
					"user_agent":   userAgent,
					"query_params": req.URL.RawQuery,
					"request_id":   c.Response().Header().Get(echo.HeaderXRequestID),
					"request_body": reqBody,
				}

				// 根据状态码决定日志级别
				log := logger.WithFields(fields)

				if statusCode >= 500 {
					log.Error("Server error")
				} else if statusCode >= 400 {
					log.Warn("Client error")
				} else {
					log.Info("Request completed")
				}

				return err
			}
		}
	}
}

// maskSensitiveData 替换敏感字段
func maskSensitiveData(data string, sensitiveFields []string) string {
	// 如果是 JSON 格式，尝试解析并替换敏感字段
	var jsonData map[string]interface{}
	if json.Unmarshal([]byte(data), &jsonData) == nil {
		// 递归处理 JSON 对象
		maskSensitiveJSON(jsonData, sensitiveFields, "")
		maskedJSON, _ := json.Marshal(jsonData)
		return string(maskedJSON)
	}

	// 如果不是 JSON 或解析失败，简单替换字符串
	result := data
	for _, field := range sensitiveFields {
		// 查找如 "password":"secret" 或 "password": "secret" 的模式
		patterns := []string{
			`"` + field + `"\s*:\s*"[^"]*"`,
			`"` + field + `"\s*:\s*'[^']*'`,
			`'` + field + `'\s*:\s*"[^"]*"`,
			`'` + field + `'\s*:\s*'[^']*'`,
		}

		for _, pattern := range patterns {
			segments := strings.Split(result, pattern)
			if len(segments) > 1 {
				// 找到匹配，组合前后片段，中间替换为敏感字段标记
				replacement := `"` + field + `":"***"`
				result = strings.Join(segments, replacement)
			}
		}
	}

	return result
}

// maskSensitiveJSON 递归替换 JSON 对象中的敏感字段
func maskSensitiveJSON(data map[string]interface{}, sensitiveFields []string, prefix string) {
	for key, value := range data {
		// 构建当前字段的完整路径
		fullPath := key
		if prefix != "" {
			fullPath = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			// 递归处理嵌套对象
			maskSensitiveJSON(v, sensitiveFields, fullPath)
		case []interface{}:
			// 处理数组
			for i, item := range v {
				if mapItem, ok := item.(map[string]interface{}); ok {
					maskSensitiveJSON(mapItem, sensitiveFields, fullPath+"."+string(rune('0'+i)))
				}
			}
		default:
			// 检查是否是敏感字段
			for _, field := range sensitiveFields {
				if strings.EqualFold(key, field) {
					// 替换敏感值
					data[key] = "***"
					break
				}
			}
		}
	}
}
