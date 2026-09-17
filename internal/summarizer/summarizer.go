package summarizer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type SummarizerClient struct {
	APIKey string
	BaseURL string
}

type SummaryRequest struct {
	Model string `json:"model"`
	Messages []Message `json:"messages"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens int `json:"max_tokens,omitempty"`
}

type Message struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

type SummaryResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type SummaryResult struct {
	Summary     string    `json:"summary"`
	KeyPoints   []string  `json:"key_points"`
	Timestamps  []TimestampPoint `json:"timestamps"`
	Duration    string    `json:"duration"`
	Tags        []string  `json:"tags"`
}

type TimestampPoint struct {
	Time string `json:"time"`
	Text string `json:"text"`
}

func NewSummarizerClient() *SummarizerClient {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	baseURL := os.Getenv("DEEPSEEK_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1"
	}
	return &SummarizerClient{
		APIKey: apiKey,
		BaseURL: baseURL,
	}
}

func (c *SummarizerClient) GetModel() string {
	if c.APIKey != "" {
		return "DeepSeek Chat"
	}
	return "Mock (no API key)"
}

func (c *SummarizerClient) Summarize(text string, options map[string]interface{}) (*SummaryResult, error) {
	if c.APIKey == "" {
		return c.generateMockSummary(text), nil
	}

	messages := []Message{
		{
			Role: "system",
			Content: "You are a professional video summarizer. Extract key points, generate concise summaries, and identify timestamps for important moments.",
		},
		{
			Role: "user",
			Content: fmt.Sprintf("Summarize this video transcript:\n\n%s\n\nPlease provide:\n1. A concise summary (2-3 sentences)\n2. Key points (bullet points)\n3. Important timestamps with descriptions\n4. Relevant tags\n\nReturn as JSON with fields: summary, key_points (array), timestamps (array with 'time' and 'text'), tags (array)", text),
		},
	}

	reqBody := SummaryRequest{
		Model: "deepseek-chat",
		Messages: messages,
		Temperature: 0.7,
		MaxTokens: 2000,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析JSON响应
	var result SummaryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return c.generateMockSummary(text), nil
	}

	if len(result.Choices) == 0 {
		return c.generateMockSummary(text), nil
	}

	content := result.Choices[0].Message.Content
	
	// 从Markdown中提取JSON
	jsonStr := extractJSON(content)
	if jsonStr == "" {
		jsonStr = content
	}

	var summary SummaryResult
	if err := json.Unmarshal([]byte(jsonStr), &summary); err != nil {
		// 如果解析失败，使用简化版本
		summary = SummaryResult{
			Summary: content,
			KeyPoints: []string{},
			Timestamps: []TimestampPoint{},
			Tags: []string{"AI", "Technology"},
		}
	}

	return &summary, nil
}

func (c *SummarizerClient) generateMockSummary(text string) *SummaryResult {
	lines := strings.Split(text, "\n")
	keyPoints := make([]string, 0)
	timestamps := make([]TimestampPoint, 0)

	for i, line := range lines {
		if len(line) > 20 {
			keyPoints = append(keyPoints, line[:min(100, len(line))])
		}
		if i%3 == 0 {
			minutes := i / 3
			timestamps = append(timestamps, TimestampPoint{
				Time: fmt.Sprintf("%d:%02d", minutes/60, minutes%60),
				Text: line[:min(50, len(line))],
			})
		}
	}

	return &SummaryResult{
		Summary: "This video covers important topics and insights.",
		KeyPoints: keyPoints[:min(5, len(keyPoints))],
		Timestamps: timestamps[:min(3, len(timestamps))],
		Duration: "10:00",
		Tags: []string{"AI", "Technology", "Tutorial"},
	}
}

func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start != -1 && end != -1 && end > start {
		return s[start:end+1]
	}
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
