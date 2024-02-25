package cmd

import (
	"errors"
)

type Format string

const (
	OldFormat Format = "old"
	NewFormat Format = "new"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *Format) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *Format) Set(v string) error {
	switch v {
	case "old", "new":
		*e = Format(v)
		return nil
	default:
		return errors.New(`must be one of "old", or "new"`)
	}
}

// Type is only used in help text
func (e *Format) Type() string {
	return "format"
}
