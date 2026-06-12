package golang

import (
	"go/ast"
	"go/parser"
	"os"
	"strings"
)

type fileScan struct {
	state      *analysisState
	file       *ast.File
	path       string
	lines      []string
	isTestFile bool
}

func (state *analysisState) processFile(path string) {
	contents, readErr := os.ReadFile(path)
	if readErr != nil {
		state.writeWarning("warning: skipping %s: %v\n", path, readErr)
		return
	}

	file, parseError := parser.ParseFile(state.fileSet, path, contents, parser.ParseComments)
	if parseError != nil {
		state.writeWarning("warning: skipping %s: %v\n", path, parseError)
		return
	}

	normalisedPath := normalisePath(path)
	state.scannedGoFiles = append(state.scannedGoFiles, normalisedPath)
	isTestFile := strings.HasSuffix(path, "_test.go")
	state.addPerFileViolations(file, normalisedPath, splitLines(contents), isTestFile)
}

func (state *analysisState) addPerFileViolations(
	file *ast.File,
	normalisedPath string,
	lines []string,
	isTestFile bool,
) {
	scan := fileScan{
		state:      state,
		file:       file,
		path:       normalisedPath,
		lines:      lines,
		isTestFile: isTestFile,
	}

	scan.addLoggingViolations()
	scan.addSecurityViolations()
	scan.addProcessViolations()
	scan.addResourceViolations()
	scan.addDataViolations()
	scan.addReturnViolations()
	scan.addParameterViolations()
	scan.addErrorViolations()
	scan.addCommentViolations()
	scan.addDomainValueViolations()
	scan.addOrderViolations()
	scan.addNamingViolations()
	scan.addTestViolations()
	scan.addFileShapeViolations()
	scan.addSpacingViolations()
}

func splitLines(contents []byte) (lines []string) {
	return strings.Split(strings.ReplaceAll(string(contents), "\r\n", "\n"), "\n")
}
