package wallarm

import "fmt"

// ExistingResourceError defines a custom error to identify when the resource has been created previously
// This error is needed to skip when we don't care about existing assets
type ExistingResourceError struct {
	Status int
	Body   string
}

// This function is required to satisfy error type
func (e *ExistingResourceError) Error() string {
	return fmt.Sprintf(
		"This resource has been created previously. The response with HTTP status code: %d and \nBody: %s",
		e.Status, e.Body)
}

// NoResourceError defines a custom error to define if no resource found.
type NoResourceError struct {
	Message string
}

// This function is required to satisfy error type.
func (e *NoResourceError) Error() string {
	return fmt.Sprintf("This resource has NOT been found. %s", e.Message)
}

// ImportAsExistsError is called when a resource exists in API and should
// be imported to work with Terraform.
func ImportAsExistsError(resourceName, id string) error {
	msg := "A resource with the ID %q already exists - to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information."
	return fmt.Errorf(msg, id, resourceName)
}
