package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/wayyoungboy/doc2mcp/internal/compiler"
	"github.com/wayyoungboy/doc2mcp/internal/index"
	"github.com/wayyoungboy/doc2mcp/internal/mcp"
)

const Version = "0.1.0"

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		help(stdout)
		return 0
	}
	switch args[0] {
	case "build":
		return runBuild(args[1:], stdout, stderr)
	case "search":
		return runSearch(args[1:], stdout, stderr)
	case "show":
		return runShow(args[1:], stdout, stderr)
	case "serve":
		return runServe(args[1:], os.Stdin, stdout, stderr)
	case "version":
		fmt.Fprintf(stdout, "Doc2MCP %s\n", Version)
		return 0
	case "help", "-h", "--help":
		help(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		help(stderr)
		return 2
	}
}

func runBuild(args []string, stdout, stderr io.Writer) int {
	parsed, err := parseOptions(args, map[string]string{"out": "", "name": ""})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if len(parsed.Positionals) != 1 {
		fmt.Fprintln(stderr, "build expects one source directory")
		return 2
	}
	out := parsed.Values["out"]
	if out == "" {
		out = "dist/" + fallbackName(parsed.Positionals[0])
	}
	if err := compiler.Build(compiler.BuildOptions{SourceDir: parsed.Positionals[0], OutDir: out, Name: parsed.Values["name"]}); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "built %s\n", out)
	return 0
}

func runSearch(args []string, stdout, stderr io.Writer) int {
	parsed, err := parseOptions(args, map[string]string{"limit": "5", "json": "false"})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if len(parsed.Positionals) < 2 {
		fmt.Fprintln(stderr, "search expects package dir and query")
		return 2
	}
	pkg, err := index.Load(parsed.Positionals[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	limit, _ := strconv.Atoi(parsed.Values["limit"])
	results := index.Search(pkg, strings.Join(parsed.Positionals[1:], " "), limit)
	if parsed.Values["json"] == "true" {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Fprintln(stdout, string(data))
		return 0
	}
	for _, result := range results {
		fmt.Fprintf(stdout, "%s\t%.2f\t%s\n", result.SectionID, result.Score, result.Snippet)
	}
	return 0
}

func runShow(args []string, stdout, stderr io.Writer) int {
	if len(args) != 2 {
		fmt.Fprintln(stderr, "show expects package dir and section id")
		return 2
	}
	pkg, err := index.Load(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	section, ok := index.SectionByID(pkg, args[1])
	if !ok {
		fmt.Fprintf(stderr, "section not found: %s\n", args[1])
		return 1
	}
	fmt.Fprintf(stdout, "# %s\n\n%s\n\nSource: %s\nChecksum: %s\n", section.Title, section.Text, section.SourcePath, section.Checksum)
	return 0
}

func runServe(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "serve expects package dir")
		return 2
	}
	pkg, err := index.Load(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := mcp.Serve(stdin, stdout, mcp.NewServer(pkg)); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

type parsedOptions struct {
	Values      map[string]string
	Positionals []string
}

func parseOptions(args []string, defaults map[string]string) (parsedOptions, error) {
	values := make(map[string]string, len(defaults))
	for key, value := range defaults {
		values[key] = value
	}
	var positionals []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "--") {
			positionals = append(positionals, arg)
			continue
		}
		nameValue := strings.TrimPrefix(arg, "--")
		name, value, inline := strings.Cut(nameValue, "=")
		if _, ok := values[name]; !ok {
			return parsedOptions{}, fmt.Errorf("unknown flag --%s", name)
		}
		if !inline {
			if i+1 >= len(args) {
				return parsedOptions{}, fmt.Errorf("missing value for --%s", name)
			}
			i++
			value = args[i]
		}
		values[name] = value
	}
	return parsedOptions{Values: values, Positionals: positionals}, nil
}

func fallbackName(path string) string {
	path = strings.TrimRight(path, "/")
	if path == "" || path == "." {
		return "docs"
	}
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func help(w io.Writer) {
	fmt.Fprintln(w, `Doc2MCP - compile docs into an MCP knowledge package

Usage:
  doc2mcp build <docs-dir> --out <package-dir> --name <name>
  doc2mcp search <package-dir> <query> [--limit 5] [--json true]
  doc2mcp show <package-dir> <section-id>
  doc2mcp serve <package-dir>
  doc2mcp version`)
}
