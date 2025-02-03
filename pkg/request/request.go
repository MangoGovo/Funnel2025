//nolint:all
package request

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

// Client 包装 Resty 客户端
type Client struct {
	*resty.Client
}

var (
	client           *Client
	clientWithoutTLS *Client
)

func newClient() *Client {
	c := &Client{
		Client: resty.New().
			SetHeader("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36").
			SetTimeout(5 * time.Second).
			SetJSONMarshaler(jsoniter.ConfigCompatibleWithStandardLibrary.Marshal).
			SetJSONUnmarshaler(jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal),
	}
	// 利用中间件实现请求日志
	c.OnAfterResponse(RestyLogMiddleware)
	return c
}

// New 初始化一个 Resty 客户端
func New() *Client {
	var once sync.Once
	once.Do(func() {
		client = newClient()
	})

	return client
}

// NewReqWithCookies 初始化一个携带 cookies 的 request 实例
func NewReqWithCookies(cookies []*http.Cookie) *resty.Request {
	return NewWithoutTLS().Request().SetCookies(cookies)
}

// NewWithoutTLS 初始化一个 Resty 客户端并跳过 TLS 证书验证
func NewWithoutTLS() *Client {
	var once sync.Once
	once.Do(func() {
		clientWithoutTLS = newClient()
		clientWithoutTLS.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	})

	return clientWithoutTLS
}

// Request 获取一个新的请求实例
func (c *Client) Request() *resty.Request {
	return c.R().EnableTrace()
}

// RestyLogMiddleware Resty日志中间件
func RestyLogMiddleware(_ *resty.Client, resp *resty.Response) error {
	if resp.IsError() {
		method := resp.Request.Method
		url := resp.Request.URL
		zap.L().Warn("请求出现错误",
			zap.String("method", method),
			zap.String("url", url),
			zap.Int64("time_spent(ms)", resp.Time().Milliseconds()),
		)
	}
	return nil
}
