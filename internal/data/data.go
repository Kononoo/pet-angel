package data

import (
	"database/sql"
	"sync"

	"pet-angel/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Data 统一数据入口（支持 database/sql 与 GORM，仓储优先使用 GORM）
// DB: 标准库 *sql.DB（部分场景可用）
// Gorm: *gorm.DB 主连接

type Data struct {
	logger *log.Helper
	DB     *sql.DB
	Gorm   *gorm.DB

	// in-memory stores（作为兜底/样例）
	userByID       map[int64]*UserDTO
	userByUsername map[string]*UserDTO
	nextUserID     int64
	mu             sync.RWMutex
}

// NewData 初始化数据库连接（GORM）
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

// UserDTO internal data model for in-memory store

type UserDTO struct {
	ID          int64
	Username    string
	Password    string
	Nickname    string
	Avatar      string
	ModelID     int64
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
