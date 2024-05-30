package validator

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type CustomRule interface {
	Tag() string
	Validate(fl validator.FieldLevel) bool
	ErrorTemplate() string
}

func NewNonZeroTimeRule() *NonZeroTimeRule {
	return &NonZeroTimeRule{}
}

type NonZeroTimeRule struct {
}

func (r *NonZeroTimeRule) Tag() string {
	return "non-zero-time"
}

func (r *NonZeroTimeRule) Validate(fl validator.FieldLevel) bool {
	switch t := fl.Field().Interface().(type) {
	case time.Time:
		return !t.IsZero()
	case *time.Time:
		return t != nil && !t.IsZero()
	default:
		return false
	}
}

func (r *NonZeroTimeRule) ErrorTemplate() string {
	return "%s is zero time"
}

func NewBeforeNowRule() *BeforeNowRule {
	return &BeforeNowRule{}
}

type BeforeNowRule struct{}

func (r *BeforeNowRule) Tag() string {
	return "before_now"
}

func (r *BeforeNowRule) Validate(fl validator.FieldLevel) bool {
	switch t := fl.Field().Interface().(type) {
	case time.Time:
		return t.Before(time.Now())
	case *time.Time:
		return t != nil && t.Before(time.Now())
	default:
		return false
	}
}

func (r *BeforeNowRule) ErrorTemplate() string {
	return "field %s must be before now"
}

func NewAfterNowRule() *AfterNowRule {
	return &AfterNowRule{}
}

type AfterNowRule struct{}

func (a AfterNowRule) Tag() string {
	return "after_now"
}

func (a AfterNowRule) Validate(fl validator.FieldLevel) bool {
	switch t := fl.Field().Interface().(type) {
	case time.Time:
		return t.After(time.Now())
	case *time.Time:
		return t != nil && t.After(time.Now())
	default:
		return false
	}
}

func (a AfterNowRule) ErrorTemplate() string {
	return "field %s must be after now"
}
