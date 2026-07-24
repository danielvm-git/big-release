// story: e08s05
package release

import (
	"errors"
	"fmt"
)

// AggregateError collects multiple errors from a lifecycle phase.
type AggregateError struct {
	errs []error
}

// NewAggregateError builds an AggregateError from one or more errors.
func NewAggregateError(errs ...error) *AggregateError {
	filtered := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			filtered = append(filtered, err)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	return &AggregateError{errs: filtered}
}

func (e *AggregateError) Error() string {
	if e == nil || len(e.errs) == 0 {
		return ""
	}
	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}
	return fmt.Sprintf("%d errors occurred: %v", len(e.errs), errors.Join(e.errs...))
}

// Unwrap returns constituent errors for errors.Is/As traversal.
func (e *AggregateError) Unwrap() []error {
	if e == nil {
		return nil
	}
	return e.errs
}
