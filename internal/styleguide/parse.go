package styleguide

import (
	"fmt"
	"os"
	"path/filepath"

	"ciphera/tools/internal/style"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

// Load parses the configured STYLE.md file under root.
func Load(root string, config Config) (document Document, err error) {
	if config.Filename == "" {
		return Document{}, fmt.Errorf("styleguide filename must not be empty")
	}

	scheme := config.IDScheme
	if err := validateIDScheme(scheme); err != nil {
		return Document{}, err
	}

	filename := config.Filename
	path := filepath.Join(root, filename)
	source, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}

	return parse(newSourceFile(filename, source), scheme)
}

// Parse parses STYLE.md source bytes.
func Parse(source []byte, config Config) (document Document, err error) {
	scheme := config.IDScheme
	if err := validateIDScheme(scheme); err != nil {
		return Document{}, err
	}

	filename := config.Filename
	if filename == "" {
		filename = defaultFilename
	}

	return parse(newSourceFile(filename, source), scheme)
}

func parse(file sourceFile, scheme style.IDScheme) (document Document, err error) {
	tree := goldmark.DefaultParser().Parse(text.NewReader(file.contents))
	events := scanMarkdown(tree, file)
	compiler := newDocumentCompiler(file, scheme)
	return compiler.compile(events)
}
