package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"ai-video-summarizer/internal/api"
	"ai-video-summarizer/internal/store"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/videos.db"
	}

	s, err := store.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer s.Close()

	r := gin.Default()
	r.Use(corsMiddleware())
	api.SetupRoutes(r, s)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("AI Video Summarizer starting on :%s", port)
	log.Fatal(r.Run(":" + port))
}
