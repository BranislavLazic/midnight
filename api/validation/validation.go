package validation

import "github.com/go-playground/validator/v10"

type Error struct {
	ErrorName string      `json:"errorName"`
	Field     string      `json:"field"`
	FieldPath string      `json:"fieldPath"`
	Value     interface{} `json:"value"`
}

func ToValidationErrors(errors validator.ValidationErrors) []Error {
	var errs []Error
	for _, e := range errors {
		errs = append(errs, Error{
			ErrorName: e.ActualTag(),
			Field:     e.StructField(),
			FieldPath: e.StructNamespace(),
			Value:     e.Value(),
		})
	}
	return errs
}
