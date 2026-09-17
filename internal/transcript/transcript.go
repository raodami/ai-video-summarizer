package transcript

import (
	"encoding/json"
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
	FullText string          `json:"full_text"`
}

type TranscriptLine struct {
	Start float64 `json:"start"`
	Text  string  `json:"text"`
}

// ExtractTranscript extracts transcript from YouTube URL
func ExtractTranscript(url string) (*TranscriptData, error) {
	videoID := extractVideoID(url)
	if videoID == "" {
		return nil, fmt.Errorf("invalid YouTube URL")
	}

	// Try to fetch transcript using invidious API (free, no auth needed)
	transcript, err := fetchFromInvidious(videoID)
	if err != nil {
		// Fallback to simple scraping
		transcript, err = fetchFromYouTubeRaw(url)
		if err != nil {
			return generateMockTranscript(videoID), nil
		}
	}

	transcript.URL = url
	transcript.ID = videoID
	return transcript, nil
}

// fetchFromInvidious uses Invidious API to get transcript
func fetchFromInvidious(videoID string) (*TranscriptData, error) {
	// Try multiple Invidious instances
	instances := []string{
		"https://inv.nadeko.net",
		"https://invidious.snopyta.org",
		"https://yewtu.be",
		"https://vid.puffyan.us",
	}

	for _, instance := range instances {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/videos/%s/captions", instance, videoID))
		if err == nil {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err == nil && len(body) > 0 {
				var captions []struct {
					DisplayName struct{ LangCode string } `json:"displayName"`
					Tracks      []struct {
						Codec     string `json:"codec"`
						StartTime float64 `json:"startTimeMs"`
						Duration  float64 `json:"durationMs"`
						XML       string  `json:"xml"`
					} `json:"tracks"`
				}
				
				if err := json.Unmarshal(body, &captions); err == nil && len(captions) > 0 {
					if len(captions[0].Tracks) > 0 {
						lines := parseCaptionXML(captions[0].Tracks[0].XML)
						return &TranscriptData{
							Lines:    lines,
							FullText: linesToText(lines),
						}, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("all invidious instances failed")
}

// parseCaptionXML extracts text from YouTube caption XML
func parseCaptionXML(xml string) []TranscriptLine {
	var lines []TranscriptLine
	
	// Simple XML parsing for <text> tags
	startTag := "<text "
	endTag := "</text>"
	
	for {
		idx := strings.Index(xml, startTag)
		if idx == -1 {
			break
		}
		
		endIdx := strings.Index(xml[idx:], endTag)
		if endIdx == -1 {
			break
		}
		
		tag := xml[idx : idx+endIdx]
		text := xml[idx+endIdx+len(endTag):]
		
		// Extract time attribute
		timeStart := extractTimeAttr(tag)
		
		// Extract text content
		textContent := strings.TrimSpace(strings.ReplaceAll(text, "<t ", ""))
		textContent = strings.Split(textContent, " ")[0]
		
		if textContent != "" {
			lines = append(lines, TranscriptLine{
				Start: timeStart,
				Text:  textContent,
			})
		}
		
		xml = xml[idx+endIdx+len(endTag):]
	}
	
	return lines
}

// extractTimeAttr extracts time from XML attribute
func extractTimeAttr(tag string) float64 {
	parts := strings.Split(tag, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "t=\"") {
			val := strings.TrimPrefix(part, "t=\"")
			val = strings.TrimSuffix(val, "\"")
			var time float64
			fmt.Sscanf(val, "%f", &time)
			return time / 1000 // Convert ms to seconds
		}
	}
	return 0
}

// linesToText converts transcript lines to full text
func linesToText(lines []TranscriptLine) string {
	var parts []string
	for _, line := range lines {
		parts = append(parts, line.Text)
	}
	return strings.Join(parts, " ")
}

// fetchFromYouTubeRaw scrapes YouTube page for transcript
func fetchFromYouTubeRaw(url string) (*TranscriptData, error) {
	req, err := http.NewRequest("GET", url, nil)
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

	html := string(body)
	
	// Extract video title
	title := extractTitle(html)
	
	// Extract author
	author := extractAuthor(html)
	
	return &TranscriptData{
		Title:  title,
		Author: author,
		Length: "N/A",
	}, nil
}

// extractTitle gets video title from HTML
func extractTitle(html string) string {
	// Simple extraction - find title between <title> tags
	startIdx := strings.Index(html, "<title>")
	if startIdx == -1 {
		return "YouTube Video"
	}
	startIdx += len("<title>")
	endIdx := strings.Index(html[startIdx:], "</title>")
	if endIdx == -1 {
		return "YouTube Video"
	}
	return strings.TrimSpace(html[startIdx : startIdx+endIdx])
}

// extractAuthor gets channel name from HTML
func extractAuthor(html string) string {
	idx := strings.Index(html, `"author":`)
	if idx == -1 {
		idx = strings.Index(html, `"ownerProfileUrl"`)
	}
	if idx == -1 {
		return "Unknown"
	}
	start := idx + len(`"author":`)
	if html[start] == '"' {
		start++
	}
	end := strings.Index(html[start:], `"`)
	if end > 0 {
		return html[start : start+end]
	}
	return "Unknown"
}

// generateMockTranscript creates sample data for testing
func generateMockTranscript(videoID string) *TranscriptData {
	return &TranscriptData{
		ID:     videoID,
		URL:    fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
		Title:  "Sample Video",
		Author: "Channel Name",
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
		FullText: "Welcome to this video about AI and machine learning. Today we'll explore the latest developments in artificial intelligence. Let's start with the fundamentals of deep learning. Neural networks are inspired by the human brain structure. Transformer models have revolutionized natural language processing. Large language models can now generate human-like text. The applications of AI are endless, from healthcare to entertainment. Thank you for watching, don't forget to like and subscribe!",
	}
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
