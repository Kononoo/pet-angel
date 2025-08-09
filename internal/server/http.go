package server

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	authv1 "pet-angel/api/auth/v1"
	avatv1 "pet-angel/api/avatar/v1"
	communityv1 "pet-angel/api/community/v1"
	greeterv1 "pet-angel/api/helloworld/v1"
	msgv1 "pet-angel/api/message/v1"
	uploadv1 "pet-angel/api/upload/v1"
	userv1 "pet-angel/api/user/v1"
	"pet-angel/internal/conf"
	"pet-angel/internal/service"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	recovery "github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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
		// 对于 proto.Message，使用 protojson.EmitUnpopulated 确保 false/0 这类零值也能输出
		if codec.Name() == "json" {
			if pm, ok := d.(proto.Message); ok {
				body, mErr := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(pm)
				if mErr != nil {
					return mErr
				}
				wrapper := struct {
					Stat   int             `json:"stat"`
					Code   int             `json:"code"`
					Msg    string          `json:"msg"`
					Reason string          `json:"reason,omitempty"`
					Data   json.RawMessage `json:"data"`
				}{Stat: 1, Code: 0, Msg: "ok", Data: body}
				if data, err = json.Marshal(&wrapper); err != nil {
					return
				}
				break
			}
		}
		reply := &Response{Stat: 1, Code: 0, Msg: "ok", Data: d}
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

// ensureDir 确保目录存在
func ensureDir(path string) error {
	if path == "" {
		return fmt.Errorf("empty path")
	}
	return os.MkdirAll(path, 0o755)
}

// saveMultipartFile 保存上传的文件
func saveMultipartFile(file multipart.File, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		return err
	}
	return nil
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, storage *conf.Storage, greeter *service.GreeterService, auth *service.AuthService, user *service.UserService, community *service.CommunityService, avatar *service.AvatarService, message *service.MessageService, upload *service.UploadService, logger log.Logger) *khttp.Server {
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
	// 静态资源服务（/static/ 前缀）
	root := "./data/assets"
	prefix := "/static/"
	if storage != nil {
		if storage.LocalRoot != "" {
			root = storage.LocalRoot
		}
		if storage.PublicPrefix != "" {
			prefix = storage.PublicPrefix
		}
	}
	_ = ensureDir(root)
	srv.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(root))))

	greeterv1.RegisterGreeterHTTPServer(srv, greeter)
	authv1.RegisterAuthServiceHTTPServer(srv, auth)
	userv1.RegisterUserServiceHTTPServer(srv, user)
	communityv1.RegisterCommunityServiceHTTPServer(srv, community)
	avatv1.RegisterAvatarServiceHTTPServer(srv, avatar)
	msgv1.RegisterMessageServiceHTTPServer(srv, message)
	uploadv1.RegisterUploadServiceHTTPServer(srv, upload)
	return srv
}
