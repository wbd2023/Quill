package golang

import (
	"go/ast"
	"go/parser"
	"strings"
)

type fileAnalysis struct {
	state      *analysisState
	file       *ast.File
	path       string
	isTestFile bool
}

func (state *analysisState) processFile(path string) {
	file, parseError := parser.ParseFile(state.fileSet, path, nil, parser.ParseComments)
	if parseError != nil {
		state.writeWarning("warning: skipping %s: %v\n", path, parseError)
		return
	}

	normalisedPath := normalisePath(path)
	state.scannedGoFiles = append(state.scannedGoFiles, normalisedPath)
	isTestFile := strings.HasSuffix(path, "_test.go")
	state.addPerFileViolations(file, normalisedPath, isTestFile)
}

func (state *analysisState) addPerFileViolations(
	file *ast.File,
	normalisedPath string,
	isTestFile bool,
) {
	analysis := fileAnalysis{
		state:      state,
		file:       file,
		path:       normalisedPath,
		isTestFile: isTestFile,
	}

	analysis.addLoggingViolations()
	analysis.addSecurityViolations()
	analysis.addProcessViolations()
	analysis.addResourceViolations()
	analysis.addDataViolations()
	analysis.addReturnViolations()
	analysis.addParameterViolations()
	analysis.addErrorViolations()
	analysis.addCommentViolations()
	analysis.addDomainIdentifierViolations()
	analysis.addOrderViolations()
	analysis.addNamingViolations()
	analysis.addTestViolations()
	analysis.addFileShapeViolations()
}
