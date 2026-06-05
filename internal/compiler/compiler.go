package compiler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wayyoungboy/doc2mcp/internal/importer"
	"github.com/wayyoungboy/doc2mcp/internal/model"
)

const BuildVersion = "0.1.0"

type BuildOptions struct {
	SourceDir string
	OutDir    string
	Name      string
}

func Build(opts BuildOptions) error {
	if opts.SourceDir == "" || opts.OutDir == "" {
		return fmt.Errorf("source and out directories are required")
	}
	name := opts.Name
	if name == "" {
		name = filepath.Base(opts.SourceDir)
	}
	docs, err := importer.ImportDirectory(opts.SourceDir)
	if err != nil {
		return err
	}
	pkg := model.Package{
		SchemaVersion: model.SchemaVersion,
		Name:          name,
		BuildVersion:  BuildVersion,
		BuiltAt:       time.Now().UTC(),
		Documents:     docs,
	}
	if err := os.MkdirAll(opts.OutDir, 0o755); err != nil {
		return err
	}
	if err := importer.CopySources(opts.SourceDir, opts.OutDir, docs); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(opts.OutDir, "doc2mcp.json"), pkg); err != nil {
		return err
	}
	return writeJSON(filepath.Join(opts.OutDir, "index.json"), flattenSections(pkg))
}

func writeJSON(path string, value interface{}) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func flattenSections(pkg model.Package) []model.Section {
	var sections []model.Section
	for _, doc := range pkg.Documents {
		sections = append(sections, doc.Sections...)
	}
	return sections
}
