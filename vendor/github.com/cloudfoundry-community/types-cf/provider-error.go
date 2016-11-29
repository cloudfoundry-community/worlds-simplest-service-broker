package cf

import (
	"fmt"
)

const (
	// ErrorServiceExists raised if instance already exists
	ErrorInstanceExists = 409

	// ErrorInstanceNotFound raised if instance not found
	ErrorInstanceNotFound = 410

	// ErrorServerException raised on server side error
	ErrorServerException = 500
)

// NewServiceProviderError factory for ServiceProviderError
func NewServiceProviderError(code int, err error) *ServiceProviderError {
	return &ServiceProviderError{Code: code, Detail: err}
}

// ServiceProviderError describes service provider error
type ServiceProviderError struct {
	Code   int
	Detail error
}

// String returns string representation of the error
func (e *ServiceProviderError) String() string {
	return fmt.Sprintf("Error: %d (%s) - %v",
		e.Code, GetServiceProviderErrorCodeName[e.Code], e.Detail.Error())
}

// GetErrorCodeName resolves error code to its string value
var GetServiceProviderErrorCodeName = map[int]string{
	409: "ErrorInstanceExists",
	410: "ErrorInstanceNotFound",
	500: "ErrorServerException",
}

// GetErrorCode resolves error name to its code
var GetServiceProviderErrorCode = map[string]int{
	"ErrorInstanceExists":   409,
	"ErrorInstanceNotFound": 410,
	"ErrorServerException":  500,
}
