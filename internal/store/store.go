package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type VideoLine struct {
	Start float64 `json:"start"`
	Text  string  `json:"text"`
}

type VideoInfo struct {
	ID        string        `json:"id"`
	URL       string        `json:"url"`
	Title     string        `json:"title"`
	Author    string        `json:"author"`
	Length    string        `json:"length"`
	Thumbnail string        `json:"thumbnail,omitempty"`
	Lines     []VideoLine   `json:"lines,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
}

type SummaryRecord struct {
	ID         string    `json:"id"`
	VideoID    string    `json:"video_id"`
	Summary    string    `json:"summary"`
	KeyPoints  string    `json:"key_points"`
	Timestamps string    `json:"timestamps"`
	Tags       string    `json:"tags"`
	Source     string    `json:"source"`
	CreatedAt  time.Time `json:"created_at"`
}

type TranslationRecord struct {
	ID           string    `json:"id"`
	VideoID      string    `json:"video_id"`
	Language     string    `json:"language"`
	Translated   string    `json:"translated"`
	CreatedAt    time.Time `json:"created_at"`
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

func NewStore(dbPath string) (*Store, error) {
	return New(dbPath)
}

func (s *Store) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS videos (
		id TEXT PRIMARY KEY,
		url TEXT NOT NULL,
		title TEXT,
		author TEXT,
		length TEXT,
		thumbnail TEXT,
		lines TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS summaries (
		id TEXT PRIMARY KEY,
		video_id TEXT NOT NULL,
		summary TEXT,
		key_points TEXT,
		timestamps TEXT,
		tags TEXT,
		source TEXT DEFAULT 'deepseek',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (video_id) REFERENCES videos(id)
	);
	
	CREATE TABLE IF NOT EXISTS translations (
		id TEXT PRIMARY KEY,
		video_id TEXT NOT NULL,
		language TEXT NOT NULL,
		translated TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (video_id) REFERENCES videos(id)
	);
	`
	_, err := s.db.Exec(schema)
	return err
}

// Video methods
func (s *Store) SaveVideo(v *VideoInfo) error {
	_, err := s.db.Exec(
		`INSERT OR REPLACE INTO videos (id, url, title, author, length, thumbnail, lines, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		v.ID, v.URL, v.Title, v.Author, v.Length, v.Thumbnail, 
		fmt.Sprintf("%v", v.Lines), v.CreatedAt,
	)
	return err
}

func (s *Store) GetVideo(id string) (*VideoInfo, error) {
	v := &VideoInfo{}
	err := s.db.QueryRow(
		`SELECT id, url, title, author, length, COALESCE(thumbnail, ''), COALESCE(lines, ''), created_at 
		 FROM videos WHERE id = ?`, id,
	).Scan(&v.ID, &v.URL, &v.Title, &v.Author, &v.Length, &v.Thumbnail, &v.Lines, &v.CreatedAt)
	return v, err
}

func (s *Store) GetAllVideos(limit int) ([]*VideoInfo, error) {
	rows, err := s.db.Query(`SELECT id, url, title, author, length, created_at FROM videos ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var videos []*VideoInfo
	for rows.Next() {
		v := &VideoInfo{}
		if err := rows.Scan(&v.ID, &v.URL, &v.Title, &v.Author, &v.Length, &v.CreatedAt); err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	return videos, nil
}

func (s *Store) GetStats() (totalVideos, totalSummaries, aiSummaries int, err error) {
	err = s.db.QueryRow("SELECT COUNT(*) FROM videos").Scan(&totalVideos)
	if err != nil {
		return 0, 0, 0, err
	}
	err = s.db.QueryRow("SELECT COUNT(*) FROM summaries").Scan(&totalSummaries)
	if err != nil {
		return 0, 0, 0, err
	}
	err = s.db.QueryRow("SELECT COUNT(*) FROM summaries WHERE source = 'deepseek'").Scan(&aiSummaries)
	return
}

// Summary methods
func (s *Store) SaveSummary(sm *SummaryRecord) error {
	_, err := s.db.Exec(
		`INSERT INTO summaries (id, video_id, summary, key_points, timestamps, tags, source, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		sm.ID, sm.VideoID, sm.Summary, sm.KeyPoints, sm.Timestamps, sm.Tags, sm.Source, sm.CreatedAt,
	)
	return err
}

func (s *Store) GetSummaries(videoID string) ([]*SummaryRecord, error) {
	rows, err := s.db.Query(`SELECT id, video_id, summary, key_points, timestamps, tags, source, created_at FROM summaries WHERE video_id = ? ORDER BY created_at DESC`, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var summaries []*SummaryRecord
	for rows.Next() {
		sm := &SummaryRecord{}
		if err := rows.Scan(&sm.ID, &sm.VideoID, &sm.Summary, &sm.KeyPoints, &sm.Timestamps, &sm.Tags, &sm.Source, &sm.CreatedAt); err != nil {
			return nil, err
		}
		summaries = append(summaries, sm)
	}
	return summaries, nil
}

// Translation methods
func (s *Store) SaveTranslation(t *TranslationRecord) error {
	_, err := s.db.Exec(
		`INSERT INTO translations (id, video_id, language, translated, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		t.ID, t.VideoID, t.Language, t.Translated, t.CreatedAt,
	)
	return err
}

func (s *Store) GetTranslations(videoID string) ([]*TranslationRecord, error) {
	rows, err := s.db.Query(`SELECT id, video_id, language, translated, created_at FROM translations WHERE video_id = ? ORDER BY created_at DESC`, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var translations []*TranslationRecord
	for rows.Next() {
		t := &TranslationRecord{}
		if err := rows.Scan(&t.ID, &t.VideoID, &t.Language, &t.Translated, &t.CreatedAt); err != nil {
			return nil, err
		}
		translations = append(translations, t)
	}
	return translations, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
