package contract

/* --------------------------------------- Make Contracts --------------------------------------- */

type MakeVariableContract struct {
	Name  string `toml:"name"`
	Value string `toml:"value"`
}

type MakeTargetContract struct {
	Name       string `toml:"name"`
	RecipeLine string `toml:"recipe_line"`
}
