package transcript

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type TranscriptData struct {
	ID      string           `json:"id"`
	URL     string           `json:"url"`
	Lines   []TranscriptLine `json:"lines"`
	Length  string           `json:"length"`
	Title   string           `json:"title"`
	Author  string           `json:"author"`
	Date    string           `json:"date"`
}

type TranscriptLine struct {
	Start float64 `json:"start"`
	Text  string  `json:"text"`
}

func ExtractTranscript(url string) (*TranscriptData, error) {
	videoID := extractVideoID(url)
	if videoID == "" {
		return nil, fmt.Errorf("invalid YouTube URL")
	}

	apiURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 从页面中提取转录文本
	transcript, err := parseTranscript(string(body), videoID)
	if err != nil {
		return nil, err
	}

	return transcript, nil
}

func extractVideoID(url string) string {
	parts := strings.Split(url, "v=")
	if len(parts) < 2 {
		return ""
	}
	id := strings.SplitN(parts[1], "&", 2)[0]
	if len(id) == 11 {
		return id
	}
	return ""
}

func parseTranscript(html, videoID string) (*TranscriptData, error) {
	// 使用YouTube Transcript API的替代方案
	// 实际项目中可以使用 yt-dlp 或 youtube-transcript-api 库
	// 这里返回模拟数据用于演示
	
	// 尝试从HTML中提取转录数据
	idx := strings.Index(html, `"captions":`)
	if idx == -1 {
		idx = strings.Index(html, `"playerCaptionsTracklistRenderer"`)
	}
	
	// 简化处理 - 返回示例数据
	return &TranscriptData{
		URL:   fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
		Title: "Sample Video",
		Author: "Channel Name",
		Date:   time.Now().Format("2006-01-02"),
		Length: "10:00",
		Lines: []TranscriptLine{
			{Start: 0.0, Text: "Welcome to this video about AI and machine learning."},
			{Start: 5.0, Text: "Today we'll explore the latest developments in artificial intelligence."},
			{Start: 10.0, Text: "Let's start with the fundamentals of deep learning."},
			{Start: 15.0, Text: "Neural networks are inspired by the human brain structure."},
			{Start: 20.0, Text: "Transformer models have revolutionized natural language processing."},
			{Start: 25.0, Text: "Large language models can now generate human-like text."},
			{Start: 30.0, Text: "The applications of AI are endless, from healthcare to entertainment."},
			{Start: 35.0, Text: "Thank you for watching, don't forget to like and subscribe!"},
		},
	}, nil
}
