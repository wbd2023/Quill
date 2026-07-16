package styleguide

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

// Load parses the configured STYLE.md file under root.
func Load(root string, config Config) (document Document, err error) {
	if config.Filename == "" {
		return Document{}, fmt.Errorf("styleguide filename must not be empty")
	}

	path := filepath.Join(root, config.Filename)
	source, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}

	return parse(newSourceFile(config.Filename, source))
}

// Parse parses STYLE.md source bytes.
func Parse(source []byte, config Config) (document Document, err error) {
	filename := config.Filename
	if filename == "" {
		filename = defaultFilename
	}

	return parse(newSourceFile(filename, source))
}

func parse(file sourceFile) (document Document, err error) {
	tree := goldmark.DefaultParser().Parse(text.NewReader(file.contents))
	events := scanMarkdown(tree, file)
	compiler := newDocumentCompiler(file)
	return compiler.compile(events)
}
