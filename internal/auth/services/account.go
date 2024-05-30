package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/krls256/dsd2024additional/internal/auth/entities"
	"github.com/krls256/dsd2024additional/internal/auth/errs"
	"github.com/krls256/dsd2024additional/internal/auth/repositories/pgsql"
	"github.com/krls256/dsd2024additional/pkg/validator"
)

func NewAccountService(accountRepository *pgsql.AccountRepository,
	validatorEngine *validator.Validator) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
		validatorEngine:   validatorEngine,
	}
}

type AccountService struct {
	accountRepository *pgsql.AccountRepository

	validatorEngine *validator.Validator
}

func (s *AccountService) GetByLogin(ctx context.Context, login string) (*entities.Account, error) {
	account, ok, err := s.accountRepository.FindBy(ctx, map[string]interface{}{
		"login": login,
	}, nil)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errs.ErrAccountNotFound
	}

	return account, nil
}

func (s *AccountService) ExitsByLogin(ctx context.Context, login string) (bool, error) {
	_, ok, err := s.accountRepository.FindBy(ctx, map[string]interface{}{
		"login": login,
	}, nil)

	return ok, err
}

func (s *AccountService) GetByID(ctx context.Context, id uuid.UUID) (*entities.Account, bool, error) {
	return s.accountRepository.FindBy(ctx, map[string]interface{}{
		"id": id,
	}, nil)
}

func (s *AccountService) All(ctx context.Context) ([]*entities.Account, error) {
	return s.accountRepository.Find(ctx, map[string]interface{}{})
}

func (s *AccountService) Create(ctx context.Context, req entities.CreateAccount) (*entities.Account, error) {
	if err := s.validatorEngine.ValidateStruct(req); err != nil {
		return nil, errors.Join(s.validatorEngine.CheckValidationPureErrors(err)...)
	}

	account, err := entities.NewAccount(req)
	if err != nil {
		return nil, err
	}

	return account, s.accountRepository.CreateNoReturn(ctx, account)
}

func (s *AccountService) Delete(ctx context.Context, req entities.ExactAccountRequest) error {
	if err := s.validatorEngine.ValidateStruct(req); err != nil {
		return errors.Join(s.validatorEngine.CheckValidationPureErrors(err)...)
	}

	return s.accountRepository.Delete(ctx, []uuid.UUID{req.ID}, map[string]interface{}{})
}

func (s *AccountService) Get(ctx context.Context, req entities.ExactAccountRequest) (*entities.Account, error) {
	acc, ok, err := s.accountRepository.FindBy(ctx, map[string]interface{}{
		"id": req.ID,
	}, nil)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errs.ErrAccountNotFound
	}

	return acc, nil
}
