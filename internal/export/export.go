package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type SummaryData struct {
	VideoID   string   `json:"video_id"`
	Title     string   `json:"title"`
	Summary   string   `json:"summary"`
	KeyPoints []string `json:"key_points"`
	Timestamps []TimestampPoint `json:"timestamps"`
	Tags      []string `json:"tags"`
	Source    string   `json:"source"`
}

type TimestampPoint struct {
	Time string `json:"time"`
	Text string `json:"text"`
}

// ExportToMarkdown converts summary to markdown format
func ExportToMarkdown(data *SummaryData) (string, error) {
	var buf bytes.Buffer
	
	buf.WriteString(fmt.Sprintf("# %s\n\n", data.Title))
	buf.WriteString(fmt.Sprintf("**Video:** %s\n\n", data.VideoID))
	buf.WriteString(fmt.Sprintf("**Source:** %s\n\n", data.Source))
	buf.WriteString(fmt.Sprintf("**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	
	buf.WriteString("---\n\n")
	
	buf.WriteString("## Summary\n\n")
	buf.WriteString(data.Summary)
	buf.WriteString("\n\n")
	
	buf.WriteString("## Key Points\n\n")
	for _, point := range data.KeyPoints {
		buf.WriteString(fmt.Sprintf("- %s\n", point))
	}
	buf.WriteString("\n")
	
	if len(data.Timestamps) > 0 {
		buf.WriteString("## Key Moments\n\n")
		for _, ts := range data.Timestamps {
			buf.WriteString(fmt.Sprintf("- **%s** %s\n", ts.Time, ts.Text))
		}
		buf.WriteString("\n")
	}
	
	if len(data.Tags) > 0 {
		buf.WriteString("## Tags\n\n")
		for _, tag := range data.Tags {
			buf.WriteString(fmt.Sprintf(" #%s", tag))
		}
		buf.WriteString("\n")
	}
	
	return buf.String(), nil
}

// ExportToJSON converts summary to JSON format
func ExportToJSON(data *SummaryData) (string, error) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// ExportToText converts summary to plain text format
func ExportToText(data *SummaryData) string {
	var buf strings.Builder
	
	buf.WriteString(fmt.Sprintf("Video: %s\n", data.Title))
	buf.WriteString(fmt.Sprintf("ID: %s\n", data.VideoID))
	buf.WriteString(fmt.Sprintf("Source: %s\n\n", data.Source))
	buf.WriteString("========================================\n\n")
	
	buf.WriteString("SUMMARY\n")
	buf.WriteString("-------\n")
	buf.WriteString(data.Summary)
	buf.WriteString("\n\n")
	
	buf.WriteString("KEY POINTS\n")
	buf.WriteString("----------\n")
	for i, point := range data.KeyPoints {
		buf.WriteString(fmt.Sprintf("%d. %s\n", i+1, point))
	}
	buf.WriteString("\n")
	
	if len(data.Timestamps) > 0 {
		buf.WriteString("KEY MOMENTS\n")
		buf.WriteString("-----------\n")
		for _, ts := range data.Timestamps {
			buf.WriteString(fmt.Sprintf("[%s] %s\n", ts.Time, ts.Text))
		}
		buf.WriteString("\n")
	}
	
	buf.WriteString("TAGS\n")
	buf.WriteString("----\n")
	for _, tag := range data.Tags {
		buf.WriteString(fmt.Sprintf("%s ", tag))
	}
	
	return buf.String()
}
