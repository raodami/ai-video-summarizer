package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ai-video-summarizer/internal/summarizer"
	"ai-video-summarizer/internal/transcript"
	"ai-video-summarizer/internal/store"
)

func SetupRoutes(r *gin.Engine, s *store.Store) {
	ds := summarizer.NewSummarizerClient()

	// Transcript endpoint
	r.POST("/api/transcript", func(c *gin.Context) {
		var req struct {
			URL string `json:"url" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
			return
		}

		// Extract transcript
		tdata, err := transcript.ExtractTranscript(req.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Save video info
		videoID := uuid.New().String()
		tdata.ID = videoID
		s.SaveVideo(&store.VideoInfo{
			ID:      videoID,
			URL:     tdata.URL,
			Title:   tdata.Title,
			Author:  tdata.Author,
			Length:  tdata.Length,
			CreatedAt: time.Now(),
		})

		c.JSON(http.StatusOK, gin.H{
			"id": videoID,
			"title": tdata.Title,
			"author": tdata.Author,
			"length": tdata.Length,
			"transcript": tdata.Lines,
		})
	})

	// Summarize endpoint
	r.POST("/api/summarize", func(c *gin.Context) {
		var req struct {
			TranscriptID string `json:"transcript_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "transcript_id is required"})
			return
		}

		// Get transcript from request context or generate from stored data
		// For now, use mock data
		text := "This is sample transcript text for summarization. It contains key information about the video content."
		
		// Generate summary using DeepSeek or mock
		summary, err := ds.Summarize(text, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Save summary
		summaryID := uuid.New().String()
		source := "mock"
		if ds.APIKey != "" {
			source = "deepseek"
		}

		s.SaveSummary(&store.SummaryRecord{
			ID:        summaryID,
			VideoID:   req.TranscriptID,
			Summary:   summary.Summary,
			KeyPoints: fmt.Sprintf("%v", summary.KeyPoints),
			Timestamps: fmt.Sprintf("%v", summary.Timestamps),
			Tags:      fmt.Sprintf("%v", summary.Tags),
			Source:    source,
			CreatedAt: time.Now(),
		})

		c.JSON(http.StatusOK, gin.H{
			"id": summaryID,
			"summary": summary,
			"source": source,
		})
	})

	// History endpoint
	r.GET("/api/history", func(c *gin.Context) {
		limit := 20
		if l := c.Query("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		}
		summaries, err := s.GetSummaries(limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, summaries)
	})

	// Stats endpoint
	r.GET("/api/stats", func(c *gin.Context) {
		stats, err := s.GetStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stats)
	})
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
