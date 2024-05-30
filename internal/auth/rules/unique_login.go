package rules

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/krls256/dsd2024additional/internal/auth/services"
)

func NewUniqueLoginRule(accountService *services.AccountService) *UniqueLoginRule {
	return &UniqueLoginRule{accountService: accountService}
}

type UniqueLoginRule struct {
	accountService *services.AccountService
}

func (r *UniqueLoginRule) Tag() string {
	return "unique_login"
}

func (r *UniqueLoginRule) Validate(fl validator.FieldLevel) bool {
	switch t := fl.Field().Interface().(type) {
	case string:
		exits, err := r.accountService.ExitsByLogin(context.Background(), t)
		if err != nil {
			return false
		}

		return !exits
	default:
		return false
	}
}

func (r *UniqueLoginRule) ErrorTemplate() string {
	return "%s field must be unique"
}
