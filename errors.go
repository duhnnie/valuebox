package valuebox

import (
	"fmt"
)

type ErrorCode string

const (
	ErrorCodeOther                = ErrorCode("other_error")
	ErrorCodeNoValueFound         = ErrorCode("no_value_found")
	ErrorCodeNonNumericArrayIndex = ErrorCode("non_numeric_array_index")
	ErrorCodeNotAMapOrSlice       = ErrorCode("not_a_map_or_slice")
)

type ResolveError struct {
	Code ErrorCode
	Path string
	Err  error
}

func (e ResolveError) Error() string {
	if e.Code == ErrorCodeOther {
		return fmt.Sprintf("(%s) %s", e.Path, e.Err)
	}

	return fmt.Sprintf("(%s) %s", e.Path, e.Code)
}

type TypeResolvingError struct {
	Type string
	Path string
}

func (e TypeResolvingError) Error() string {
	return fmt.Sprintf("(%s) can't resolve to type: %s", e.Path, e.Type)
}
