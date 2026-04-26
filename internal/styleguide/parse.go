package styleguide

import (
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func compileDocument(
	contents []byte,
	styleGuide Config,
) (document Document, err error) {
	root := goldmark.DefaultParser().Parse(text.NewReader(contents))
	walkState := newDocumentWalkState(styleGuide.RequirementIDFormat)
	if err = gast.Walk(root, walkState.walk(contents)); err != nil {
		return Document{}, err
	}

	return walkState.finish()
}
