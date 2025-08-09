package server

import (
	"net/http"
	"strconv"
	"strings"

	authv1 "pet-angel/api/auth/v1"
	avatv1 "pet-angel/api/avatar/v1"
	communityv1 "pet-angel/api/community/v1"
	greeterv1 "pet-angel/api/helloworld/v1"
	userv1 "pet-angel/api/user/v1"
	"pet-angel/internal/conf"
	"pet-angel/internal/service"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	recovery "github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
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
func ResponseEncoder(w http.ResponseWriter, r *http.Request, d interface{}) (err error) {
	codec, _ := khttp.CodecForRequest(r, "Accept")
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
func ErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	codec, _ := khttp.CodecForRequest(r, "Accept")
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

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, auth *service.AuthService, user *service.UserService, community *service.CommunityService, avatar *service.AvatarService, logger log.Logger) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
		),
		// CORS 过滤器：允许跨域与预检
		khttp.Filter(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				next.ServeHTTP(w, r)
			})
		}),
		khttp.ResponseEncoder(ResponseEncoder),
		khttp.ErrorEncoder(ErrorEncoder),
	}
	if c.Http.Network != "" {
		opts = append(opts, khttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, khttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := khttp.NewServer(opts...)
	greeterv1.RegisterGreeterHTTPServer(srv, greeter)
	authv1.RegisterAuthServiceHTTPServer(srv, auth)
	userv1.RegisterUserServiceHTTPServer(srv, user)
	communityv1.RegisterCommunityServiceHTTPServer(srv, community)
	avatv1.RegisterAvatarServiceHTTPServer(srv, avatar)
	return srv
}
