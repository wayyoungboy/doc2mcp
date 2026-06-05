package importer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/wayyoungboy/doc2mcp/internal/model"
)

var headingPattern = regexp.MustCompile(`(?m)^(#{1,6})\s+(.+?)\s*$`)

func ImportDirectory(root string) ([]model.Document, error) {
	var paths []string
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".md" || ext == ".markdown" || ext == ".txt" {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Strings(paths)

	docs := make([]model.Document, 0, len(paths))
	for _, path := range paths {
		doc, err := importFile(root, path)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func importFile(root, path string) (model.Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return model.Document{}, err
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return model.Document{}, err
	}
	rel = filepath.ToSlash(rel)
	text := string(data)
	sections := splitSections(rel, text)
	title := rel
	if len(sections) > 0 && sections[0].Title != "" {
		title = sections[0].Title
	}
	doc := model.Document{
		ID:         rel,
		SourcePath: rel,
		Title:      title,
		Checksum:   checksum(data),
		Sections:   sections,
	}
	return doc, nil
}

func splitSections(sourcePath, text string) []model.Section {
	matches := headingPattern.FindAllStringSubmatchIndex(text, -1)
	if len(matches) == 0 || !strings.HasSuffix(strings.ToLower(sourcePath), ".md") {
		return []model.Section{makeSection(sourcePath, sourcePath, "Document", text, 0, len(text))}
	}
	sections := make([]model.Section, 0, len(matches))
	for i, match := range matches {
		titleStart, titleEnd := match[4], match[5]
		contentStart := match[1]
		contentEnd := len(text)
		if i+1 < len(matches) {
			contentEnd = matches[i+1][0]
		}
		title := strings.TrimSpace(text[titleStart:titleEnd])
		content := strings.TrimSpace(text[contentStart:contentEnd])
		sections = append(sections, makeSection(sourcePath, sourcePath+"#"+slug(title), title, content, contentStart, contentEnd))
	}
	return sections
}

func makeSection(sourcePath, id, title, text string, start, end int) model.Section {
	return model.Section{
		ID:         id,
		DocID:      sourcePath,
		SourcePath: sourcePath,
		Title:      title,
		Text:       strings.TrimSpace(text),
		StartByte:  start,
		EndByte:    end,
		Checksum:   checksum([]byte(text)),
	}
}

func slug(value string) string {
	lower := strings.ToLower(value)
	var out strings.Builder
	lastDash := false
	for _, r := range lower {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			out.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			out.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(out.String(), "-")
}

func checksum(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func CopySources(root, out string, docs []model.Document) error {
	for _, doc := range docs {
		from := filepath.Join(root, filepath.FromSlash(doc.SourcePath))
		to := filepath.Join(out, "sources", filepath.FromSlash(doc.SourcePath))
		if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
			return err
		}
		data, err := os.ReadFile(from)
		if err != nil {
			return fmt.Errorf("read source %s: %w", from, err)
		}
		if err := os.WriteFile(to, data, 0o644); err != nil {
			return fmt.Errorf("write source %s: %w", to, err)
		}
	}
	return nil
}
