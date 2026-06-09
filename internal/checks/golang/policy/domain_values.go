package policy

// DomainValueConfig defines domain value conversion policy.
type DomainValueConfig struct {
	RequiredConstructors DomainValueConstructors
}

// DomainValueConstructors maps domain value types to approved constructors.
type DomainValueConstructors map[string][]string
