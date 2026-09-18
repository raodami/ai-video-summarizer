package export

import (
	"fmt"
	"strings"
	"time"
)

type ExportFormat string

const (
	FormatMD    ExportFormat = "markdown"
	FormatJSON  ExportFormat = "json"
	FormatTXT   ExportFormat = "text"
	FormatSRT   ExportFormat = "srt"
	FormatVTT   ExportFormat = "vtt"
)

type ExportResult struct {
	Content string    `json:"content"`
	MimeType string   `json:"mime_type"`
	Filename string   `json:"filename"`
	Format   ExportFormat `json:"format"`
}

func ExportSummary(title, summary, keyPoints, timestamps, format string) (*ExportResult, error) {
	switch ExportFormat(format) {
	case FormatMD:
		return exportMarkdown(title, summary, keyPoints, timestamps)
	case FormatJSON:
		return exportJSON(title, summary, keyPoints, timestamps)
	case FormatTXT:
		return exportText(title, summary, keyPoints, timestamps)
	case FormatSRT:
		return exportSRT(title, summary, timestamps)
	case FormatVTT:
		return exportVTT(title, summary, timestamps)
	default:
		return exportMarkdown(title, summary, keyPoints, timestamps)
	}
}

func exportMarkdown(title, summary, keyPoints, timestamps string) (*ExportResult, error) {
	content := fmt.Sprintf("# %s\n\n", title)
	content += fmt.Sprintf("## Summary\n\n%s\n\n", summary)
	
	if keyPoints != "" {
		content += "## Key Points\n\n"
		for _, point := range strings.Split(keyPoints, "\n") {
			point = strings.TrimSpace(point)
			if point != "" {
				content += fmt.Sprintf("- %s\n", point)
			}
		}
		content += "\n"
	}
	
	if timestamps != "" {
		content += "## Timestamps\n\n"
		for _, line := range strings.Split(timestamps, "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				content += fmt.Sprintf("%s\n", line)
			}
		}
	}
	
	return &ExportResult{
		Content:  content,
		MimeType: "text/markdown",
		Filename: fmt.Sprintf("%s_%s.md", sanitize(title), time.Now().Format("2006-01-02")),
		Format:   FormatMD,
	}, nil
}

func exportJSON(title, summary, keyPoints, timestamps string) (*ExportResult, error) {
	content := fmt.Sprintf(`{
  "title": "%s",
  "summary": "%s",
  "key_points": [%s],
  "timestamps": "%s",
  "exported_at": "%s"
}`, sanitizeJSON(title), sanitizeJSON(summary), sanitizeJSON(keyPoints), sanitizeJSON(timestamps), time.Now().Format(time.RFC3339))
	
	return &ExportResult{
		Content:  content,
		MimeType: "application/json",
		Filename: fmt.Sprintf("%s_%s.json", sanitize(title), time.Now().Format("2006-01-02")),
		Format:   FormatJSON,
	}, nil
}

func exportText(title, summary, keyPoints, timestamps string) (*ExportResult, error) {
	content := fmt.Sprintf("%s\n\n=== SUMMARY ===\n\n%s", title, summary)
	
	if keyPoints != "" {
		content += "\n\n=== KEY POINTS ===\n\n" + keyPoints
	}
	
	if timestamps != "" {
		content += "\n\n=== TIMESTAMPS ===\n\n" + timestamps
	}
	
	return &ExportResult{
		Content:  content,
		MimeType: "text/plain",
		Filename: fmt.Sprintf("%s_%s.txt", sanitize(title), time.Now().Format("2006-01-02")),
		Format:   FormatTXT,
	}, nil
}

func exportSRT(title, summary, timestamps string) (*ExportResult, error) {
	var sb strings.Builder
	sb.WriteString("WEBVTT\n\n")
	sb.WriteString(fmt.Sprintf("NOTE\n%s\n\n", title))
	
	if timestamps != "" {
		lines := strings.Split(timestamps, "\n")
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			idx := i + 1
			sb.WriteString(fmt.Sprintf("%d\n00:00:0%d,000 --> 00:00:%d,000\n%s\n\n", idx, idx, idx+1, line))
		}
	}
	
	return &ExportResult{
		Content:  sb.String(),
		MimeType: "text/vtt",
		Filename: fmt.Sprintf("%s_%s.vtt", sanitize(title), time.Now().Format("2006-01-02")),
		Format:   FormatVTT,
	}, nil
}

func exportVTT(title, summary, timestamps string) (*ExportResult, error) {
	return exportSRT(title, summary, timestamps)
}

func sanitize(s string) string {
	return strings.NewReplacer("/", "-", "\\", "-", ":", "-", " ", "_").Replace(s)
}

func sanitizeJSON(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
