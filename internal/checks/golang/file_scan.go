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

	normalised := normalisePath(path)
	state.scannedGoFiles = append(state.scannedGoFiles, normalised)
	isTestFile := strings.HasSuffix(path, "_test.go")
	state.addPerFileViolations(file, normalised, splitLines(contents), isTestFile)
}

func (state *analysisState) addPerFileViolations(
	file *ast.File,
	normalised string,
	lines []string,
	isTestFile bool,
) {
	scan := fileScan{
		state:      state,
		file:       file,
		path:       normalised,
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
