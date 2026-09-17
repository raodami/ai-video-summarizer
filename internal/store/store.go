package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type VideoInfo struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Length      string    `json:"length"`
	Thumbnail   string    `json:"thumbnail,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type SummaryRecord struct {
	ID          string    `json:"id"`
	VideoID     string    `json:"video_id"`
	Summary     string    `json:"summary"`
	KeyPoints   string    `json:"key_points"`
	Timestamps  string    `json:"timestamps"`
	Tags        string    `json:"tags"`
	Source      string    `json:"source"` // "deepseek" or "mock"
	CreatedAt   time.Time `json:"created_at"`
}

type Store struct {
	db *sql.DB
}

func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &Store{db: db}
	if err := store.initSchema(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *Store) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS videos (
		id TEXT PRIMARY KEY,
		url TEXT NOT NULL,
		title TEXT,
		author TEXT,
		length TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS summaries (
		id TEXT PRIMARY KEY,
		video_id TEXT NOT NULL,
		summary TEXT,
		key_points TEXT,
		timestamps TEXT,
		tags TEXT,
		source TEXT DEFAULT 'mock',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (video_id) REFERENCES videos(id)
	);
	CREATE INDEX IF NOT EXISTS idx_summaries_video ON summaries(video_id);
	`
	_, err := s.db.Exec(schema)
	return err
}

func (s *Store) SaveVideo(v *VideoInfo) error {
	_, err := s.db.Exec(
		"INSERT OR REPLACE INTO videos (id, url, title, author, length) VALUES (?, ?, ?, ?, ?)",
		v.ID, v.URL, v.Title, v.Author, v.Length,
	)
	return err
}

func (s *Store) SaveSummary(smry *SummaryRecord) error {
	_, err := s.db.Exec(
		"INSERT INTO summaries (id, video_id, summary, key_points, timestamps, tags, source) VALUES (?, ?, ?, ?, ?, ?, ?)",
		smry.ID, smry.VideoID, smry.Summary, smry.KeyPoints, smry.Timestamps, smry.Tags, smry.Source,
	)
	return err
}

func (s *Store) GetSummary(id string) (*SummaryRecord, error) {
	var smry SummaryRecord
	err := s.db.QueryRow(
		"SELECT id, video_id, summary, key_points, timestamps, tags, source, created_at FROM summaries WHERE id = ?",
		id,
	).Scan(&smry.ID, &smry.VideoID, &smry.Summary, &smry.KeyPoints, &smry.Timestamps, &smry.Tags, &smry.Source, &smry.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &smry, nil
}

func (s *Store) GetSummaries(limit int) ([]*SummaryRecord, error) {
	rows, err := s.db.Query("SELECT id, video_id, summary, key_points, timestamps, tags, source, created_at FROM summaries ORDER BY created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []*SummaryRecord
	for rows.Next() {
		var s SummaryRecord
		if err := rows.Scan(&s.ID, &s.VideoID, &s.Summary, &s.KeyPoints, &s.Timestamps, &s.Tags, &s.Source, &s.CreatedAt); err != nil {
			continue
		}
		summaries = append(summaries, &s)
	}
	return summaries, rows.Err()
}

func (s *Store) GetVideo(url string) (*VideoInfo, error) {
	var v VideoInfo
	err := s.db.QueryRow("SELECT id, url, title, author, length, created_at FROM videos WHERE url = ?", url).
		Scan(&v.ID, &v.URL, &v.Title, &v.Author, &v.Length, &v.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *Store) GetStats() (map[string]interface{}, error) {
	var videoCount int
	var summaryCount int
	var deepseekCount int
	
	s.db.QueryRow("SELECT COUNT(*) FROM videos").Scan(&videoCount)
	s.db.QueryRow("SELECT COUNT(*) FROM summaries").Scan(&summaryCount)
	s.db.QueryRow("SELECT COUNT(*) FROM summaries WHERE source = 'deepseek'").Scan(&deepseekCount)

	return map[string]interface{}{
		"total_videos":   videoCount,
		"total_summaries": summaryCount,
		"ai_summaries":   deepseekCount,
	}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
