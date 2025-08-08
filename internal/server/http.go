package server

import (
	stdhttp "net/http"
	"strconv"
	"strings"

	v1 "pet-angel/api/helloworld/v1"
	"pet-angel/internal/conf"
	"pet-angel/internal/service"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// 统一响应结构
type Response struct {
	Stat   int         `json:"stat"`             // 状态：1-成功，0-失败
	Code   int         `json:"code"`             // HTTP 状态码
	Msg    string      `json:"msg"`              // 消息
	Reason string      `json:"reason,omitempty"` // 错误原因
	Data   interface{} `json:"data"`             // 业务数据
}

// ContentType returns the content-type with base prefix.
func ContentType(subtype string) string {
	return strings.Join([]string{"application", subtype}, "/")
}

// IsCallback 判断是否为回调接口
func IsCallback(url string) bool {
	return strings.HasSuffix(url, "/callback")
}

// ResponseEncoder 统一响应编码器
func ResponseEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, d interface{}) (err error) {
	codec, _ := http.CodecForRequest(r, "Accept")
	w.Header().Set("Content-Type", ContentType(codec.Name()))
	var data []byte

	switch {
	case IsCallback(r.RequestURI):
		// 回调接口直接返回原始数据
		if data, err = codec.Marshal(d); err != nil {
			return
		}
	default:
		// 普通接口包装成统一响应格式
		reply := &Response{
			Stat: 1,
			Code: 0,
			Msg:  "ok",
			Data: d,
		}
		if data, err = codec.Marshal(reply); err != nil {
			return
		}
	}

	_, err = w.Write(data)
	return
}

// ErrorEncoder 错误响应编码器
func ErrorEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
	codec, _ := http.CodecForRequest(r, "Accept")
	w.Header().Set("Content-Type", ContentType(codec.Name()))

	se := errors.FromError(err)
	code := int(se.Code)
	var body []byte

	switch {
	case IsCallback(r.RequestURI):
		// 回调接口返回简单错误信息
		response := strings.Replace(se.Message, "{code}", strconv.Itoa(code), 1)
		response = strings.Replace(response, "{reason}", se.Reason, 1)
		body = []byte(response)
	default:
		// 普通接口返回统一错误格式
		reply := &Response{
			Stat:   0,
			Code:   code,
			Msg:    se.Message,
			Reason: se.Reason,
			Data:   nil,
		}
		if data, marshalErr := codec.Marshal(reply); marshalErr != nil {
			body = []byte(`{"stat":0,"code":500,"msg":"Internal Server Error","data":null}`)
		} else {
			body = data
		}
	}

	w.WriteHeader(stdhttp.StatusOK)
	w.Write(body)
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
		http.ResponseEncoder(ResponseEncoder),
		http.ErrorEncoder(ErrorEncoder),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
