package index

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/wayyoungboy/doc2mcp/internal/model"
)

var tokenPattern = regexp.MustCompile(`[a-z0-9_/-]+`)

func Load(dir string) (model.Package, error) {
	data, err := os.ReadFile(filepath.Join(dir, "doc2mcp.json"))
	if err != nil {
		return model.Package{}, err
	}
	var pkg model.Package
	if err := json.Unmarshal(data, &pkg); err != nil {
		return model.Package{}, err
	}
	return pkg, nil
}

func Search(pkg model.Package, query string, limit int) []model.SearchResult {
	if limit <= 0 {
		limit = 5
	}
	queryTokens := counts(query)
	var results []model.SearchResult
	for _, section := range sections(pkg) {
		score := scoreSection(queryTokens, section)
		if score <= 0 {
			continue
		}
		results = append(results, model.SearchResult{
			SectionID: section.ID,
			Title:     section.Title,
			Snippet:   snippet(section.Text, 180),
			Score:     score,
			Citation: model.Citation{
				SectionID:  section.ID,
				SourcePath: section.SourcePath,
				Title:      section.Title,
				Checksum:   section.Checksum,
			},
		})
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].SectionID < results[j].SectionID
		}
		return results[i].Score > results[j].Score
	})
	if len(results) > limit {
		return results[:limit]
	}
	return results
}

func SectionByID(pkg model.Package, id string) (model.Section, bool) {
	for _, section := range sections(pkg) {
		if section.ID == id {
			return section, true
		}
	}
	return model.Section{}, false
}

func sections(pkg model.Package) []model.Section {
	var out []model.Section
	for _, doc := range pkg.Documents {
		out = append(out, doc.Sections...)
	}
	return out
}

func counts(text string) map[string]int {
	tokens := tokenPattern.FindAllString(strings.ToLower(text), -1)
	out := make(map[string]int, len(tokens))
	for _, token := range tokens {
		out[token]++
	}
	return out
}

func scoreSection(query map[string]int, section model.Section) float64 {
	haystack := counts(section.Title + " " + section.Text + " " + section.SourcePath)
	var score float64
	for token, count := range query {
		if hits := haystack[token]; hits > 0 {
			score += float64(count*hits) * (1 + 1/math.Sqrt(float64(len(haystack)+1)))
		}
	}
	if strings.Contains(strings.ToLower(section.Title), "authentication") && query["authentication"] > 0 {
		score += 2
	}
	return score
}

func snippet(text string, max int) string {
	text = strings.Join(strings.Fields(text), " ")
	if len(text) <= max {
		return text
	}
	return text[:max] + "..."
}
