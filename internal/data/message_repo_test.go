package data

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestGorm builds an in-memory sqlite gorm DB with required tables.
func setupTestGorm(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	// Minimal schema for messages and users
	if err := gdb.Exec(`
CREATE TABLE users (id INTEGER PRIMARY KEY, coins INTEGER DEFAULT 0);
CREATE TABLE messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  sender INTEGER NOT NULL,
  message_type INTEGER NOT NULL,
  is_locked INTEGER NOT NULL,
  unlock_coins INTEGER NOT NULL,
  content TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE user_unlock_records (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  message_id INTEGER NOT NULL,
  coins_spent INTEGER NOT NULL
);
`).Error; err != nil {
		t.Fatal(err)
	}
	return gdb
}

func TestUnlockMessage(t *testing.T) {
	gdb := setupTestGorm(t)
	sqlDB, _ := gdb.DB()
	d := &Data{Gorm: gdb, DB: sqlDB}
	repo := NewMessageRepo(d)
	// seed user 1 with 100 coins and a locked note cost 20
	if err := gdb.Exec(`INSERT INTO users(id,coins) VALUES (1,100);`).Error; err != nil {
		t.Fatal(err)
	}
	if err := gdb.Exec(`INSERT INTO messages(user_id,sender,message_type,is_locked,unlock_coins,content) VALUES (1,1,1,1,20,'hi');`).Error; err != nil {
		t.Fatal(err)
	}
	// find message id
	var mid int64
	_ = gdb.Raw(`SELECT id FROM messages LIMIT 1`).Scan(&mid).Error
	remain, msg, err := repo.UnlockMessage(context.Background(), 1, mid)
	if err != nil {
		t.Fatalf("unlock err: %v", err)
	}
	if remain != 80 {
		t.Fatalf("remain want 80 got %d", remain)
	}
	if msg == nil || msg.IsLocked {
		t.Fatal("message should be unlocked")
	}
	// repeat unlock should be idempotent
	remain2, _, err := repo.UnlockMessage(context.Background(), 1, mid)
	if err != nil {
		t.Fatalf("second unlock err: %v", err)
	}
	if remain2 != 80 {
		t.Fatalf("remain should stay 80 got %d", remain2)
	}
}
