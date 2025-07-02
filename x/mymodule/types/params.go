package types

// NewParams creates a new Params instance
func NewParams(mymoduleDisabled bool) Params {
	return Params{
		MymoduleDisabled: mymoduleDisabled,
	}
}

const (
	defaultModuleDisabled = false
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(defaultModuleDisabled)
}

// Validate validates the set of params
func (p Params) ValidateBasic() error {
	return nil
}
