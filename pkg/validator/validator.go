package validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"reflect"
	"strings"
	"sync"
)

var ErrValidation = errors.New("")

type Validator struct {
	*validator.Validate
	once          sync.Once
	errorMessages map[string]string
}

type TaggedError struct {
	Tag string
	Err error
}

func New(rules ...CustomRule) (*Validator, error) {
	v := &Validator{
		Validate:      validator.New(),
		errorMessages: defaultErrors(),
	}

	for _, rule := range rules {
		if err := v.register(rule.Tag(), rule.Validate); err != nil {
			return nil, err
		}

		v.errorMessages[rule.Tag()] = rule.ErrorTemplate()
	}

	v.init()

	return v, nil
}

func (v *Validator) AddCustomRules(rules ...CustomRule) error {
	for _, rule := range rules {
		if err := v.register(rule.Tag(), rule.Validate); err != nil {
			return err
		}

		v.errorMessages[rule.Tag()] = rule.ErrorTemplate()
	}

	return nil
}

const NameSplitterCount = 2

func (v *Validator) init() {
	v.once.Do(func() {
		v.SetTagName("validate")
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", NameSplitterCount)[0]

			if name == "-" {
				return ""
			}

			return name
		})
	})
}

func (v *Validator) register(tag string, fn validator.Func) error {
	return v.Validate.RegisterValidation(tag, fn)
}

func (v *Validator) ValidateStruct(obj interface{}) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.Struct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(sliceValidateError, 0)

		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}

		if len(validateRet) == 0 {
			return nil
		}

		return validateRet
	case reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32, reflect.Float64,
		reflect.Func, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Interface,
		reflect.Invalid, reflect.Map, reflect.String, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
		fallthrough
	default:
		return nil
	}
}

const (
	ZeroWC = 0
	OneWC  = 1
	TwoWC  = 2
)

func (v *Validator) CheckValidationPureErrors(err error) []error {
	return lo.Map(v.CheckValidationErrors(err), func(item TaggedError, index int) error {
		return item.Err
	})
}

func (v *Validator) CheckValidationErrors(err error) (e []TaggedError) {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		e = append(e, TaggedError{Tag: InvalidTag, Err: err})
	}

	errs, ok := err.(validator.ValidationErrors)

	if !ok {
		e = append(e, TaggedError{Tag: InvalidTag, Err: err})

		return e
	}

	for _, validationError := range errs {
		message := v.errorMessages[validationError.Tag()]
		formattedMessage := ""

		switch strings.Count(message, "%s") {
		case ZeroWC:
			formattedMessage = message
		case OneWC:
			formattedMessage = fmt.Sprintf(message, validationError.Field())
		case TwoWC:
			formattedMessage = fmt.Sprintf(message, validationError.Field(), validationError.Param())
		}

		e = append(e, TaggedError{Tag: validationError.Tag(), Err: fmt.Errorf("%w%s", ErrValidation, formattedMessage)})
	}

	return e
}
