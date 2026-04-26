package policy

type ControlPlaneConfig struct {
	QualityFile       string
	VariableContracts []MakeVariableContract
	TargetContracts   []MakeTargetContract
}

type MakeVariableContract struct {
	Name  string
	Value string
}

type MakeTargetContract struct {
	Name       string
	RecipeLine string
}
