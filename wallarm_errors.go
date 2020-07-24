package wallarm

import "fmt"

// ExistingResourceError defines a custom error to identify when the resource has been created previously.
// This error is needed to skip when we don't care about existing assets.
type ExistingResourceError struct {
	Status int
	Body   string
}

// This function is required to satisfy a custom error type.
func (e *ExistingResourceError) Error() string {
	return fmt.Sprintf(
		"This resource has been created previously. The response with HTTP status code: %d and \nBody: %s",
		e.Status, e.Body)
}

// NoResourceError is raised when no resource found.
type NoResourceError struct {
	Message string
}

// This function is required to satisfy a custom error type.
func (e *NoResourceError) Error() string {
	return fmt.Sprintf("This resource has NOT been found. %s", e.Message)
}
