package contract

import "fmt"

// RequireFlag returns a VALIDATION_ERROR if value is empty.
func RequireFlag(name, value string) error {
	if value == "" {
		return NewError(ValidationError, fmt.Sprintf("--%s is required", name))
	}
	return nil
}
