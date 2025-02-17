package db

import (
	"context"
	"fmt"
	"log"
	"os"

	redis "github.com/go-redis/redis/v8" // 使用具体版本
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Rdb *redis.Client
var ctx = context.Background()

func Initialize() error {
	// Database initialization
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" {
		return fmt.Errorf("database configuration missing")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return err
	}

	// Redis initialization
	redisAddr := os.Getenv("REDIS_HOST") + ":6379"
	if redisAddr == ":6379" {
		return fmt.Errorf("redis host not defined")
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err = Rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return nil
}
