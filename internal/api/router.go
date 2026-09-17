package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ai-video-summarizer/internal/export"
	"ai-video-summarizer/internal/summarizer"
	"ai-video-summarizer/internal/transcript"
	"ai-video-summarizer/internal/store"
)

func SetupRoutes(r *gin.Engine, s *store.Store, ds *summarizer.SummarizerClient) {
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
			ID:        videoID,
			URL:       tdata.URL,
			Title:     tdata.Title,
			Author:    tdata.Author,
			Length:    tdata.Length,
			CreatedAt: time.Now(),
		})

		c.JSON(http.StatusOK, gin.H{
			"id":         videoID,
			"title":      tdata.Title,
			"author":     tdata.Author,
			"length":     tdata.Length,
			"transcript": tdata.Lines,
			"full_text":  tdata.FullText,
		})
	})

	// Summarize endpoint
	r.POST("/api/summarize", func(c *gin.Context) {
		var req struct {
			TranscriptID string `json:"transcript_id" binding:"required"`
			FullText     string `json:"full_text"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "transcript_id is required"})
			return
		}

		// Get full text from request or generate from stored data
		text := req.FullText
		if text == "" {
			text = "This is sample transcript text for summarization."
		}

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
			ID:         summaryID,
			VideoID:    req.TranscriptID,
			Summary:    summary.Summary,
			KeyPoints:  toJSON(summary.KeyPoints),
			Timestamps: toJSON(summary.Timestamps),
			Tags:       toJSON(summary.Tags),
			Source:     source,
			CreatedAt:  time.Now(),
		})

		c.JSON(http.StatusOK, gin.H{
			"id":       summaryID,
			"summary":  summary,
			"source":   source,
			"model":    ds.GetModel(),
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

	// Export endpoints
	r.GET("/api/export/:id/markdown", func(c *gin.Context) {
		id := c.Param("id")
		smry, err := s.GetSummary(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "summary not found"})
			return
		}
		data, _ := parseSummaryRecord(smry)
		markdown, err := export.ExportToMarkdown(data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "text/markdown")
		c.String(http.StatusOK, markdown)
	})

	r.GET("/api/export/:id/json", func(c *gin.Context) {
		id := c.Param("id")
		smry, err := s.GetSummary(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "summary not found"})
			return
		}
		data, _ := parseSummaryRecord(smry)
		json, err := export.ExportToJSON(data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, json)
	})

	r.GET("/api/export/:id/text", func(c *gin.Context) {
		id := c.Param("id")
		smry, err := s.GetSummary(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "summary not found"})
			return
		}
		data, _ := parseSummaryRecord(smry)
		text := export.ExportToText(data)
		c.String(http.StatusOK, text)
	})
}

func toJSON(v interface{}) string {
	if v == nil {
		return "[]"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func parseSummaryRecord(smry *store.SummaryRecord) (*export.SummaryData, error) {
	var keyPoints []string
	var timestamps []export.TimestampPoint
	var tags []string

	json.Unmarshal([]byte(smry.KeyPoints), &keyPoints)
	json.Unmarshal([]byte(smry.Timestamps), &timestamps)
	json.Unmarshal([]byte(smry.Tags), &tags)

	return &export.SummaryData{
		VideoID:  smry.VideoID,
		Summary:  smry.Summary,
		KeyPoints: keyPoints,
		Timestamps: timestamps,
		Tags:     tags,
		Source:   smry.Source,
	}, nil
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
