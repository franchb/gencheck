package benchmark

import (
	"errors"

	"github.com/abice/gencheck"
)

func (s SingleString) Validate() error {

	var vErrors gencheck.ValidationErrors

	// BEGIN Entry Validations
	// required

	if s.Entry == "" {
		return append(vErrors, gencheck.NewFieldError("SingleString", "Entry", "required", errors.New("is required")))
	}

	// END Entry Validations

	if len(vErrors) > 0 {
		return vErrors
	}
	return nil
}

func (s TestString) Validate() error {

	var vErrors gencheck.ValidationErrors

	// BEGIN Required Validations
	// required

	if s.Required == "" {
		return append(vErrors, gencheck.NewFieldError("TestString", "Required", "required", errors.New("is required")))
	}

	// END Required Validations

	// BEGIN Len Validations
	// len

	if !(len(s.Len) == 10) {
		return append(vErrors, gencheck.NewFieldError("TestString", "Len", "len", errors.New("length mismatch")))
	}

	// END Len Validations

	// BEGIN Min Validations
	// min

	if len(s.Min) < 5 {
		return append(vErrors, gencheck.NewFieldError("TestString", "Min", "min", errors.New("length was less than 5")))
	}

	// END Min Validations

	// BEGIN Max Validations
	// max

	if len(s.Max) > 100 {
		return append(vErrors, gencheck.NewFieldError("TestString", "Max", "max", errors.New("length was more than 100")))
	}

	// END Max Validations

	if len(vErrors) > 0 {
		return vErrors
	}
	return nil
}