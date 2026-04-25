package styleguide

import (
	"ciphera/tools/internal/profile"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func compileDocument(
	contents []byte,
	styleGuide profile.StyleGuideConfig,
) (document documentModel, err error) {
	root := goldmark.DefaultParser().Parse(text.NewReader(contents))
	walkState := newDocumentWalkState(styleGuide.RequirementIDFormat)
	if err = gast.Walk(root, walkState.walk(contents)); err != nil {
		return documentModel{}, err
	}

	return walkState.finish()
}
