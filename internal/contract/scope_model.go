package contract

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	LevelRequired       Level = "required"
	LevelRecommendation Level = "recommendation"
)

const (
	CheckProfileRequired CheckProfile = "required"
	CheckProfileAll      CheckProfile = "all"
)

const (
	ScopeUnknown Scope = ""
	ScopeApp     Scope = "app"
	ScopeTools   Scope = "tools"
	ScopeAll     Scope = "all"
)

/* -------------------------------------------- Types ------------------------------------------- */

type Level string

type CheckProfile string

type Scope string

/* --------------------------------------- Scope Coverage --------------------------------------- */

func ScopeCovers(requested Scope, required Scope) (allowed bool) {
	switch required {
	case ScopeAll:
		return true

	case ScopeApp:
		return requested == ScopeApp || requested == ScopeAll

	case ScopeTools:
		return requested == ScopeTools || requested == ScopeAll

	default:
		return false
	}
}
