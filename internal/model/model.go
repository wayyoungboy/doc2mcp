package model

import "time"

const SchemaVersion = "v1"

type Package struct {
	SchemaVersion string     `json:"schema_version"`
	Name          string     `json:"name"`
	BuildVersion  string     `json:"build_version"`
	BuiltAt       time.Time  `json:"built_at"`
	Documents     []Document `json:"documents"`
}

type Document struct {
	ID         string    `json:"id"`
	SourcePath string    `json:"source_path"`
	Title      string    `json:"title"`
	Checksum   string    `json:"checksum"`
	Sections   []Section `json:"sections"`
}

type Section struct {
	ID         string `json:"id"`
	DocID      string `json:"doc_id"`
	SourcePath string `json:"source_path"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	StartByte  int    `json:"start_byte"`
	EndByte    int    `json:"end_byte"`
	Checksum   string `json:"checksum"`
}

type Citation struct {
	SectionID  string `json:"section_id"`
	SourcePath string `json:"source_path"`
	Title      string `json:"title"`
	Checksum   string `json:"checksum"`
}

type SearchResult struct {
	SectionID string   `json:"section_id"`
	Title     string   `json:"title"`
	Snippet   string   `json:"snippet"`
	Score     float64  `json:"score"`
	Citation  Citation `json:"citation"`
}
