package batch

import (
	"context"
	"log"
	"strconv"
	"sync"

	"advertising-system/db"
	"advertising-system/models"

	"gorm.io/gorm"
)

func ProcessCounts() {
	ctx := context.Background()

	keys, err := db.Rdb.Keys(ctx, "ad:*:counts").Result()
	if err != nil {
		log.Fatalf("Failed to get keys from Redis: %v", err)
		return
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20)

	for _, key := range keys {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(key string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			processKey(ctx, key)

		}(key)
	}

	wg.Wait()
}

func processKey(ctx context.Context, key string) {

	adIDStr := key[3 : len(key)-7]
	adID, err := strconv.ParseInt(adIDStr, 10, 64)
	if err != nil {
		log.Printf("Failed to parse ad ID from key %s: %v", key, err)
		return
	}

	fields, err := db.Rdb.HGetAll(ctx, key).Result()
	if err != nil {
		log.Printf("Error retrieving count data from Redis for key %s: %v", key, err)
		return
	}

	cpcCount, _ := strconv.Atoi(fields["cpc_count"])
	cpmCount, _ := strconv.Atoi(fields["cpm_count"])

	if err := updateDatabase(adID, cpcCount, cpmCount); err != nil {
		log.Printf("Failed to update database for ad %d: %v", adID, err)
		return
	}

	// 删除Redis键
	if err := db.Rdb.Del(ctx, key).Err(); err != nil {
		log.Printf("Failed to delete key %s from Redis: %v", key, err)
	}
}

func updateDatabase(adID int64, cpcCount, cpmCount int) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	result := tx.Model(&models.Advertisement{}).
		Where("id = ?", adID).
		Updates(map[string]interface{}{
			"cpc_count": gorm.Expr("cpc_count + ?", cpcCount),
			"cpm_count": gorm.Expr("cpm_count + ?", cpmCount),
		})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit() // 提交事务
	return nil
}
