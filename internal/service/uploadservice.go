package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	uploadv1 "pet-angel/api/upload/v1"
	"pet-angel/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// UploadService 实现本地文件上传

type UploadService struct {
	uploadv1.UnimplementedUploadServiceServer
	storage *conf.Storage
	logger  *log.Helper
}

func NewUploadService(storage *conf.Storage, l log.Logger) *UploadService {
	return &UploadService{storage: storage, logger: log.NewHelper(l)}
}

// ensureDir 确保目录存在
func ensureDir(path string) error { return os.MkdirAll(path, 0o755) }

// saveMultipartFile 保存上传的文件
func saveMultipartFile(file multipart.File, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	return err
}

// ------- 校验配置 -------
const (
	maxImageMB = 50
	maxVideoMB = 1000
)

var (
	allowedImageMIMEs = map[string]bool{
		"image/jpeg": true, "image/png": true, "image/gif": true, "image/webp": true,
	}
	allowedVideoMIMEs = map[string]bool{
		"video/mp4": true, "video/quicktime": true,
	}
)

func mimeToExt(ct string) string {
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	case "video/quicktime":
		return ".mov"
	default:
		return ""
	}
}

func (s *UploadService) UploadFile(ctx context.Context, in *uploadv1.UploadFileRequest) (*uploadv1.UploadFileReply, error) {
	root := "./data/assets"
	prefix := "/static/"
	if s.storage != nil {
		if s.storage.LocalRoot != "" {
			root = s.storage.LocalRoot
		}
		if s.storage.PublicPrefix != "" {
			prefix = s.storage.PublicPrefix
		}
	}
	_ = ensureDir(root)

	ts, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil, http.ErrNoCookie
	}
	ht, ok := ts.(*khttp.Transport)
	if !ok {
		return nil, http.ErrNoCookie
	}
	req := ht.Request()
	// 解析 multipart 文件
	if err := req.ParseMultipartForm(32 << 20); err != nil { // 32MB
		return nil, err
	}
	file, header, err := req.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	category := req.FormValue("type")
	if category == "" {
		category = strings.ToLower(strings.TrimSpace(in.GetType()))
	}
	if category == "" {
		category = "image"
	}
	// 归一化：avatar 视作 image 类目保存
	isAvatar := category == "avatar"
	isVideo := category == "video"
	if isAvatar {
		category = "avatar"
	} else if isVideo {
		category = "video"
	} else {
		category = "image"
	}

	// sniff content-type
	sniff := make([]byte, 512)
	n, _ := file.Read(sniff)
	_, _ = file.Seek(0, io.SeekStart)
	contentType := http.DetectContentType(sniff[:n])

	// 校验 MIME 与大小
	switch category {
	case "image", "avatar":
		if !allowedImageMIMEs[contentType] {
			return nil, errors.New("unsupported image type")
		}
		if header.Size > 0 && header.Size > int64(maxImageMB)*1024*1024 {
			return nil, errors.New("image too large")
		}
	case "video":
		if !allowedVideoMIMEs[contentType] {
			return nil, errors.New("unsupported video type")
		}
		if header.Size > 0 && header.Size > int64(maxVideoMB)*1024*1024 {
			return nil, errors.New("video too large")
		}
	}
	subdir := filepath.Join(category, time.Now().Format("2006/01/02"))
	baseDir := filepath.Join(root, subdir)
	if err := ensureDir(baseDir); err != nil {
		return nil, err
	}
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		if e := mimeToExt(contentType); e != "" {
			ext = e
		} else {
			ext = ".bin"
		}
	}
	filename := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), os.Getpid(), ext)
	dst := filepath.Join(baseDir, filename)
	if err := saveMultipartFile(file, dst); err != nil {
		return nil, err
	}
	url := strings.TrimRight(prefix, "/") + "/" + filepath.ToSlash(filepath.Join(subdir, filename))
	return &uploadv1.UploadFileReply{Url: url}, nil
}

func (s *UploadService) GetPresign(ctx context.Context, in *uploadv1.GetPresignRequest) (*uploadv1.GetPresignReply, error) {
	// 目前本地上传，不提供预签；返回占位
	return &uploadv1.GetPresignReply{Url: "", Method: "", Headers: map[string]string{}, FinalUrl: ""}, nil
}

func (s *UploadService) UploadDone(ctx context.Context, in *uploadv1.UploadDoneRequest) (*uploadv1.UploadDoneReply, error) {
	// 直传登记占位；直接回显
	return &uploadv1.UploadDoneReply{Url: in.GetUrl()}, nil
}
