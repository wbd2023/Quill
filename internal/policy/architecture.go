package policy

type ArchitectureConfig struct {
	Layers []ArchitectureLayer
}

type ArchitectureLayer struct {
	Name         string
	PackageRoots []string
	MayImport    []string
}
