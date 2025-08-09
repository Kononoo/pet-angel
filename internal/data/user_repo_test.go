package data

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserGorm(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := gdb.Exec(`
CREATE TABLE user_follows (follower_id INTEGER NOT NULL, followee_id INTEGER NOT NULL, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY(follower_id,followee_id));
`).Error; err != nil {
		t.Fatal(err)
	}
	return gdb
}

func TestFollowIdempotent(t *testing.T) {
	gdb := setupUserGorm(t)
	d := &Data{Gorm: gdb}
	repo := NewUserRepo(d)
	if err := repo.Follow(context.Background(), 1, 2); err != nil {
		t.Fatal(err)
	}
	if err := repo.Follow(context.Background(), 1, 2); err != nil {
		t.Fatal(err)
	}
	var cnt int64
	if err := gdb.Table("user_follows").Count(&cnt).Error; err != nil {
		t.Fatal(err)
	}
	if cnt != 1 {
		t.Fatalf("expect 1 row, got %d", cnt)
	}
}
