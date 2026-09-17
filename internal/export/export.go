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

// ExportToSRT converts transcript to SRT subtitle format
func ExportToSRT(lines []TranscriptLine) string {
	var buf strings.Builder
	for i, line := range lines {
		start := formatSRTTime(line.Start)
		end := formatSRTTime(line.Start + 3.0)
		buf.WriteString(fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, start, end, line.Text))
	}
	return buf.String()
}

// ExportToVTT converts transcript to WebVTT subtitle format
func ExportToVTT(lines []TranscriptLine) string {
	var buf strings.Builder
	buf.WriteString("WEBVTT\n\n")
	for i, line := range lines {
		start := formatVTTTime(line.Start)
		end := formatVTTTime(line.Start + 3.0)
		buf.WriteString(fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, start, end, line.Text))
	}
	return buf.String()
}

func formatSRTTime(seconds float64) string {
	hours := int(seconds / 3600)
	minutes := int((int(seconds) % 3600) / 60)
	secs := int(seconds) % 60
	millis := int((seconds - float64(int(seconds))) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, secs, millis)
}

func formatVTTTime(seconds float64) string {
	hours := int(seconds / 3600)
	minutes := int((int(seconds) % 3600) / 60)
	secs := int(seconds) % 60
	millis := int((seconds - float64(int(seconds))) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, secs, millis)
}

type TranscriptLine struct {
	Start float64 `json:"start"`
	Text  string  `json:"text"`
}
