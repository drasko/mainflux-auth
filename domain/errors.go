package domain

import "fmt"

type AuthError struct {
	Code    int
	Message string
}

func (e *AuthError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}
