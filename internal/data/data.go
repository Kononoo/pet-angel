package data

import (
	"context"
	"database/sql"
	"sync"

	"pet-angel/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Data 统一数据入口（支持 database/sql 与 GORM，仓储优先使用 GORM）
// DB: 标准库 *sql.DB（部分场景可用）
// Gorm: *gorm.DB 主连接
// Minio: MinIO 对象存储客户端

type Data struct {
	logger *log.Helper
	DB     *sql.DB
	Gorm   *gorm.DB

	Minio       *minio.Client
	MinioBucket string

	// in-memory stores（作为兜底/样例）
	userByID       map[int64]*UserDTO
	userByUsername map[string]*UserDTO
	nextUserID     int64
	mu             sync.RWMutex
}

// NewData 初始化数据库与 MinIO 连接
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(logger)
	var db *sql.DB
	var gdb *gorm.DB
	var err error
	if c != nil && c.Database != nil && c.Database.Source != "" {
		dsn := c.Database.Source
		gdb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, nil, err
		}
		db, err = gdb.DB()
		if err != nil {
			return nil, nil, err
		}
		if err = db.Ping(); err != nil {
			return nil, nil, err
		}
	}

	d := &Data{
		logger:         l,
		DB:             db,
		Gorm:           gdb,
		userByID:       make(map[int64]*UserDTO),
		userByUsername: make(map[string]*UserDTO),
		nextUserID:     1,
	}

	cleanup := func() {
		l.Info("closing the data resources")
		if d.DB != nil {
			_ = d.DB.Close()
		}
	}
	return d, cleanup, nil
}

// InitMinio 创建 MinIO 客户端并确保桶存在
func (d *Data) InitMinio(ctx context.Context, mc *conf.Minio) error {
	if mc == nil || mc.Endpoint == "" {
		return nil
	}
	cli, err := minio.New(mc.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(mc.AccessKey, mc.SecretKey, ""),
		Secure: mc.UseSsl,
	})
	if err != nil {
		return err
	}
	d.Minio = cli
	d.MinioBucket = mc.Bucket
	// 确保 bucket 存在
	exists, err := cli.BucketExists(ctx, mc.Bucket)
	if err != nil {
		return err
	}
	if !exists {
		if err := cli.MakeBucket(ctx, mc.Bucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
	}
	return nil
}

// UserDTO internal data model for in-memory store

type UserDTO struct {
	ID          int64
	Username    string
	Password    string
	Nickname    string
	Avatar      string
	ModelID     int64
	ModelURL    string
	PetName     string
	PetAvatar   string
	PetSex      int32
	Kind        string
	Weight      int32
	Hobby       string
	Description string
	Coins       int32
	CreatedAt   int64 // unix seconds
}
