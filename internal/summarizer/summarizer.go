package summarizer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type SupportedLanguage struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
}

var Languages = []SupportedLanguage{
	{Code: "en", Name: "English"},
	{Code: "zh", Name: "Chinese"},
	{Code: "ja", Name: "Japanese"},
	{Code: "ko", Name: "Korean"},
	{Code: "es", Name: "Spanish"},
	{Code: "fr", Name: "French"},
	{Code: "de", Name: "German"},
	{Code: "pt", Name: "Portuguese"},
	{Code: "ru", Name: "Russian"},
	{Code: "ar", Name: "Arabic"},
	{Code: "hi", Name: "Hindi"},
	{Code: "th", Name: "Thai"},
	{Code: "vi", Name: "Vietnamese"},
	{Code: "id", Name: "Indonesian"},
	{Code: "ms", Name: "Malay"},
	{Code: "tl", Name: "Filipino"},
}

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
	Source      string    `json:"source"`
}

type TimestampPoint struct {
	Time string `json:"time"`
	Text string `json:"text"`
}

func NewClient() *SummarizerClient {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	return &SummarizerClient{
		APIKey: apiKey,
		BaseURL: "https://api.deepseek.com/v1",
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
		return generateMockSummary(text), nil
	}

	prompt := buildPrompt(text, options)

	reqBody, _ := json.Marshal(SummaryRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{"system", "You are a helpful assistant that summarizes video transcripts."},
			{"user", prompt},
		},
		Temperature: 0.3,
		MaxTokens: 2048,
	})

	resp, err := http.Post(
		fmt.Sprintf("%s/chat/completions", c.BaseURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return generateMockSummary(text), nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result SummaryResponse
	json.Unmarshal(body, &result)

	if len(result.Choices) > 0 {
		return parseSummaryResponse(result.Choices[0].Message.Content), nil
	}

	return generateMockSummary(text), nil
}

func buildPrompt(text string, options map[string]interface{}) string {
	language := "en"
	if lang, ok := options["language"].(string); ok && lang != "" {
		language = lang
	}

	instructions := map[string]string{
		"en": "Please summarize this video transcript. Include: 1) A brief summary 2) Key points 3) Important timestamps 4) Relevant tags",
		"zh": "请总结这段视频字幕。包括：1) 简要概述 2) 关键点 3) 重要时间戳 4) 相关标签",
		"ja": "この動画の要約を作成してください。以下の要素を含めてください：1) まとめ 2) 重要なポイント 3) 重要なタイムスタンプ 4) 関連タグ",
		"ko": "이 영상 요약을 만들어 주세요. 다음 요소를 포함하세요: 1) 요약 2) 주요 포인트 3) 중요한 타임스탬프 4) 관련 태그",
		"es": "Resume este transcripción de video. Incluye: 1) Un breve resumen 2) Puntos clave 3) Marcas de tiempo importantes 4) Etiquetas relevantes",
		"fr": "Résumez cette transcription vidéo. Incluez: 1) Un bref résumé 2) Points clés 3) Horodatages importants 4) Tags pertinents",
		"de": "Fassen Sie dieses Video-Zusammenfassung zusammen. Beinhaltet: 1) Kurze Zusammenfassung 2) Schlüsselinformationen 3) Wichtige Zeitstempel 4) Relevante Tags",
		"pt": "Resuma esta transcrição de vídeo. Inclua: 1) Breve resumo 2) Pontos principais 3) Marcações de tempo importantes 4) Tags relevantes",
		"ru": "Подведите итог этой расшифровки видео. Включите: 1) Краткое изложение 2) Ключевые моменты 3) Важные временные метки 4) Соответствующие теги",
		"ar": "لخص هذا النص من الفيديو. يتضمن: 1) ملخص موجز 2) النقاط الرئيسية 3) الطوابع الزمنية المهمة 4) الوسوم ذات الصلة",
		"hi": "इस वीडियो ट्रांसक्रिप्ट का सारांश दें। शामिल करें: 1) संक्षिप्त सारांश 2) मुख्य बिंदु 3) महत्वपूर्ण समय स्तंभ 4) प्रासंगिक टैग",
	}

	inst, ok := instructions[language]
	if !ok {
		inst = instructions["en"]
	}

	return fmt.Sprintf("%s\n\n%s", inst, text)
}

func parseSummaryResponse(content string) *SummaryResult {
	result := &SummaryResult{
		Summary: content,
		Source:  "deepseek",
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
			result.KeyPoints = append(result.KeyPoints, strings.TrimPrefix(strings.TrimPrefix(line, "-"), "*"))
		} else if strings.Contains(line, ":") && len(line) < 50 {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				result.Timestamps = append(result.Timestamps, TimestampPoint{
					Time: strings.TrimSpace(parts[0]),
					Text: strings.TrimSpace(parts[1]),
				})
			}
		}
		if strings.HasPrefix(line, "#") {
			result.Tags = append(result.Tags, strings.TrimPrefix(line, "#"))
		}
	}

	return result
}

func generateMockSummary(text string) *SummaryResult {
	return &SummaryResult{
		Summary: "This video covers important topics about the subject matter. The speaker provides detailed insights and practical examples.",
		KeyPoints: []string{
			"Introduction to the main topic",
			"Key concepts and explanations",
			"Practical examples and case studies",
			"Summary and conclusions",
		},
		Timestamps: []TimestampPoint{
			{Time: "00:00", Text: "Introduction"},
			{Time: "02:30", Text: "Main topic overview"},
			{Time: "05:45", Text: "Detailed explanation"},
			{Time: "08:20", Text: "Examples and demonstration"},
			{Time: "10:00", Text: "Conclusion"},
		},
		Tags: []string{"education", "tutorial", "technology"},
		Source: "mock",
	}
}
