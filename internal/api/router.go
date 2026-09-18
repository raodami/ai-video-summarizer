package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ai-video-summarizer/internal/export"
	"ai-video-summarizer/internal/queue"
	"ai-video-summarizer/internal/summarizer"
	"ai-video-summarizer/internal/transcript"
	"ai-video-summarizer/internal/store"
	"ai-video-summarizer/internal/ws"
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

		// Send WebSocket update
		ws.ManagerInstance.SendProgress(videoID, 50, "success", "Transcript extracted")

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
		ws.ManagerInstance.SendProgress(req.TranscriptID, 70, "processing", "Generating AI summary...")
		summary, err := ds.Summarize(text, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Save summary
		summaryID := uuid.New().String()
		s.SaveSummary(&store.SummaryRecord{
			ID:        summaryID,
			VideoID:   req.TranscriptID,
			Summary:   summary.Summary,
			KeyPoints: strings.Join(summary.KeyPoints, "\n"),
			Timestamps: "",
			Tags:      strings.Join(summary.Tags, ","),
			Source:    summary.Source,
			CreatedAt: time.Now(),
		})

		// Update progress
		ws.ManagerInstance.SendProgress(req.TranscriptID, 100, "success", "Summary generated")

		c.JSON(http.StatusOK, gin.H{
			"id":          summaryID,
			"video_id":    req.TranscriptID,
			"summary":     summary.Summary,
			"key_points":  summary.KeyPoints,
			"timestamps":  "",
			"tags":        summary.Tags,
			"source":      summary.Source,
		})
	})

	// Batch queue endpoints
	r.POST("/api/batch/add", func(c *gin.Context) {
		var req struct {
			URLs []string `json:"urls" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "URLs are required"})
			return
		}

		bq := queue.GetInstance()
		tasks := make([]*queue.BatchTask, 0, len(req.URLs))
		for _, url := range req.URLs {
			task := bq.AddTask(url)
			tasks = append(tasks, task)
		}
		c.JSON(http.StatusOK, gin.H{"tasks": tasks, "count": len(tasks)})
	})

	r.GET("/api/batch/tasks", func(c *gin.Context) {
		bq := queue.GetInstance()
		tasks := bq.GetAllTasks()
		c.JSON(http.StatusOK, gin.H{"tasks": tasks, "count": len(tasks)})
	})

	r.DELETE("/api/batch/clear", func(c *gin.Context) {
		bq := queue.GetInstance()
		cleared := bq.ClearCompleted()
		c.JSON(http.StatusOK, gin.H{"cleared": cleared})
	})

	// Export endpoint
	r.POST("/api/export", func(c *gin.Context) {
		var req struct {
			VideoID  string `json:"video_id" binding:"required"`
			Format   string `json:"format"`
			Title    string `json:"title"`
			Summary  string `json:"summary"`
			KeyPoints string `json:"key_points"`
			Timestamps string `json:"timestamps"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Required fields missing"})
			return
		}
		if req.Format == "" {
			req.Format = "markdown"
		}
		if req.Title == "" {
			req.Title = "Video Summary"
		}

		result, err := export.ExportSummary(req.Title, req.Summary, req.KeyPoints, req.Timestamps, req.Format)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	// Stats endpoint
	r.GET("/api/stats", func(c *gin.Context) {
		totalVideos, totalSummaries, aiSummaries, err := s.GetStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"total_videos":    totalVideos,
			"total_summaries": totalSummaries,
			"ai_summaries":    aiSummaries,
		})
	})

	// History endpoint
	r.GET("/api/history", func(c *gin.Context) {
		limit := 20
		if l := c.Query("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		}
		videos, err := s.GetAllVideos(limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, videos)
	})

	// Languages endpoint
	r.GET("/api/languages", func(c *gin.Context) {
		languages := []map[string]string{
			{"code": "en", "name": "English"},
			{"code": "es", "name": "Spanish"},
			{"code": "fr", "name": "French"},
			{"code": "de", "name": "German"},
			{"code": "zh", "name": "Chinese"},
			{"code": "ja", "name": "Japanese"},
			{"code": "ko", "name": "Korean"},
			{"code": "ar", "name": "Arabic"},
			{"code": "pt", "name": "Portuguese"},
			{"code": "ru", "name": "Russian"},
			{"code": "hi", "name": "Hindi"},
			{"code": "it", "name": "Italian"},
			{"code": "nl", "name": "Dutch"},
			{"code": "tr", "name": "Turkish"},
			{"code": "vi", "name": "Vietnamese"},
			{"code": "th", "name": "Thai"},
		}
		c.JSON(http.StatusOK, languages)
	})

	// WS endpoint
	ws.SetupWS(r)
}

func generateID() string {
	return uuid.New().String()
}
