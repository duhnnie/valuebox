package valuebox

import "fmt"

type Error string

const (
	ErrorResolveInvalidParams     = Error("second argument needs to be a string or a slice of strings.")
	ErrorResolveInvalidFirstParam = Error("\"target\" parameter should be a \"map[string]interface{}\" when second argument is not an empty string nor empty string slice")
	ErrorTargetIsNIL              = Error("\"target\" parameter is nil")
)

func (err Error) Error() string {
	return string(err)
}

type ErrorNoValueFound string

func (e ErrorNoValueFound) Error() string {
	return fmt.Sprintf("no value was found for \"%s\"", string(e))
}

type ErrorCantResolveToType struct {
	Type string
	Name string
}

func (e ErrorCantResolveToType) Error() string {
	return fmt.Sprintf("can't resolve \"%s\" to type: %s", e.Name, e.Type)
}
