package http

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client HTTP 客户端封装
type Client struct {
	client *resty.Client
}

// ClientOption 客户端配置选项
type ClientOption struct {
	Timeout       time.Duration
	RetryCount    int
	RetryWaitTime time.Duration
}

// DefaultOption 默认配置
var DefaultOption = ClientOption{
	Timeout:       time.Second * 10,
	RetryCount:    3,
	RetryWaitTime: time.Second,
}

// NewClient 创建新的 HTTP 客户端
func NewClient(opt *ClientOption) *Client {
	if opt == nil {
		opt = &DefaultOption
	}

	client := resty.New().
		SetTimeout(opt.Timeout).
		SetRetryCount(opt.RetryCount).
		SetRetryWaitTime(opt.RetryWaitTime)

	return &Client{client: client}
}

// 添加包级别的默认客户端
var defaultClient = NewClient(nil)

// 设置默认客户端的选项
func SetDefaultOption(opt *ClientOption) {
	defaultClient = NewClient(opt)
}

// Post 发送 POST 请求
func (c *Client) Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*resty.Response, error) {
	req := c.client.R().
		SetContext(ctx).
		SetBody(body).
		SetHeader("Content-Type", "application/json")

	if headers != nil {
		req.SetHeaders(headers)
	}

	return req.Post(url)
}

// Put 发送 PUT 请求
func (c *Client) Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*resty.Response, error) {
	req := c.client.R().
		SetContext(ctx).
		SetBody(body).
		SetHeader("Content-Type", "application/json")

	if headers != nil {
		req.SetHeaders(headers)
	}

	return req.Put(url)
}

// Get 发送 GET 请求
func (c *Client) Get(ctx context.Context, url string, params map[string]string) (*resty.Response, error) {
	return c.client.R().
		SetContext(ctx).
		SetQueryParams(params).
		Get(url)
}

// Delete 发送 DELETE 请求
func (c *Client) Delete(ctx context.Context, url string) (*resty.Response, error) {
	return c.client.R().
		SetContext(ctx).
		Delete(url)
}

// SetHeader 设置请求头
func (c *Client) SetHeader(key, value string) *Client {
	c.client.SetHeader(key, value)
	return c
}

// SetHeaders 批量设置请求头
func (c *Client) SetHeaders(headers map[string]string) *Client {
	c.client.SetHeaders(headers)
	return c
}

// SetAuthToken 设置认证 Token
func (c *Client) SetAuthToken(token string) *Client {
	c.client.SetAuthToken(token)
	return c
}

// Response 统一响应结构
type Response struct {
	Body       string
	StatusCode int
}

// 新增包级别的快捷方法（带统一响应）
// 添加处理响应的公共方法
func handleResponse(resp *resty.Response, err error) (*Response, error) {
    if err != nil {
        return nil, err
    }

    if !resp.IsSuccess() {
        return nil, fmt.Errorf("请求失败: status=%d, body=%s", resp.StatusCode(), resp.String())
    }

    return &Response{
        Body:       resp.String(),
        StatusCode: resp.StatusCode(),
    }, nil
}

// 修改包级别的方法
func Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error) {
    return handleResponse(defaultClient.Post(ctx, url, body, headers))
}

func Get(ctx context.Context, url string, params map[string]string) (*Response, error) {
    return handleResponse(defaultClient.Get(ctx, url, params))
}

func Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error) {
    return handleResponse(defaultClient.Put(ctx, url, body, headers))
}

func Delete(ctx context.Context, url string) (*Response, error) {
    return handleResponse(defaultClient.Delete(ctx, url))
}
