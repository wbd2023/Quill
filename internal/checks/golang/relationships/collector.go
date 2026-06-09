package relationships

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/checks/golang/analysis"
)

type Collector struct {
	pathClassifier         analysis.PathClassifier
	interfaces             map[string]interfaceDeclaration
	mocks                  map[string][]methodDeclaration
	implementations        map[string][]methodDeclaration
	implementationBindings []implementationBinding
}

func NewCollector(pathClassifier analysis.PathClassifier) (collector *Collector) {
	return &Collector{
		pathClassifier:         pathClassifier,
		interfaces:             make(map[string]interfaceDeclaration),
		mocks:                  make(map[string][]methodDeclaration),
		implementations:        make(map[string][]methodDeclaration),
		implementationBindings: make([]implementationBinding, 0),
	}
}

func (collector *Collector) Collect(fileSet *token.FileSet, file *ast.File, path string) {
	collectInterfaces(fileSet, file, path, collector.pathClassifier, collector.interfaces)
	collectMockMethods(fileSet, file, path, collector.pathClassifier, collector.mocks)
	collectImplementationMethods(
		fileSet,
		file,
		path,
		collector.pathClassifier,
		collector.implementations,
	)
	collectImplementationBindings(
		fileSet,
		file,
		path,
		collector.pathClassifier,
		&collector.implementationBindings,
	)
}

func (collector *Collector) Violations() (violations []analysis.Violation) {
	violations = append(
		violations,
		checkMockOrderAgainstInterfaces(collector.interfaces, collector.mocks)...,
	)
	violations = append(
		violations,
		checkImplementationOrderAgainstInterfaces(
			collector.interfaces,
			collector.implementations,
			collector.implementationBindings,
		)...,
	)

	return violations
}
