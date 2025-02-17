package main

import (
	"log"
	"net/http"
	"os"

	"advertising-system/api"
	"advertising-system/cmd/batch"
	"advertising-system/db"

	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
)

func main() {
	if err := db.Initialize(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/advertisers/{id}", api.GetAdvertiser).Methods("GET")
	r.HandleFunc("/api/advertisements/{id}", api.GetAdvertisement).Methods("GET")
	r.HandleFunc("/api/advertisements/{id}/impression", api.RecordImpression).Methods("GET")
	r.HandleFunc("/api/advertisements/{id}/click", api.RecordClick).Methods("GET")
	r.HandleFunc("/api/ad/request", api.RequestAd).Methods("GET")

	// 设置cron任务
	c := cron.New()
	_, err := c.AddFunc("@every 10s", batch.ProcessCounts)
	if err != nil {
		log.Fatalf("Error scheduling batch process: %v", err)
	}
	// 启动cron调度器
	c.Start()
	defer c.Stop() // 确保在main函数退出时停止调度器

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
