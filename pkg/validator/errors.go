package validator

import (
	"fmt"
	"strings"
)

func defaultErrors() map[string]string {
	return map[string]string{
		"required":      "field %s is required",
		"email_custom":  "email %s is not valid",
		"str_gt":        "field %s must have greater than %s characters",
		"str_lt":        "field %s must have less than %s characters",
		"has_lowercase": "field %s must have at least one small character",
		"has_uppercase": "field %s must have at least one big character",
		"has_special":   "field %s must have at least one special character",
		"oneof":         "field %s must have value one of allowed list: %s",
		"gte":           "field %s must be greater or equal than %s",
		"gt":            "field %s must be greater than %s",
		"lte":           "field %s must be less or equal than %s",
		"lt":            "field %s must be less than %s",
		"url":           "field %s must be an url",
		"uuid":          "field %s must be an uuid",
		"min":           "min len for field %s: %s",
		"max":           "max len for field %s: %s",
	}
}

type sliceValidateError []error

func (err sliceValidateError) Error() string {
	errMsgs := []string{}

	for i, e := range err {
		if e == nil {
			continue
		}

		errMsgs = append(errMsgs, fmt.Sprintf("[%d]: %s", i, e.Error()))
	}

	return strings.Join(errMsgs, "\n")
}
