package domain

import "fmt"

type ServiceError struct {
	Code    int
	Message string
}

func (e *ServiceError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}
