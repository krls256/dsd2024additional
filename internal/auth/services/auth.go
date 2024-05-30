package services

import (
	"context"
	"errors"
	"github.com/krls256/dsd2024additional/internal/auth/entities"
	"github.com/krls256/dsd2024additional/internal/auth/errs"
	"github.com/krls256/dsd2024additional/pkg/auth"
	pkgRepositories "github.com/krls256/dsd2024additional/pkg/repositories/pgsql"
	"github.com/krls256/dsd2024additional/pkg/validator"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

func NewAuthService(
	tokenRepository *pkgRepositories.BaseRepository[*entities.Token],
	accountService *AccountService,
	sessionService *SessionService,
	jwtAuthorizer *auth.JWTAuthorizer,
	validatorEngine *validator.Validator) *AuthService {
	return &AuthService{
		accountService:  accountService,
		tokenRepository: tokenRepository,
		sessionService:  sessionService,

		jwtAuthorizer:   jwtAuthorizer,
		validatorEngine: validatorEngine,
	}
}

type AuthService struct {
	tokenRepository *pkgRepositories.BaseRepository[*entities.Token]
	accountService  *AccountService
	sessionService  *SessionService

	jwtAuthorizer   *auth.JWTAuthorizer
	validatorEngine *validator.Validator
}

func (s *AuthService) Login(ctx context.Context, req entities.LoginRequest) (*auth.JWT, error) {
	if err := s.validatorEngine.ValidateStruct(req); err != nil {
		return nil, errors.Join(s.validatorEngine.CheckValidationPureErrors(err)...)
	}

	account, err := s.accountService.GetByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}

	if !entities.CheckPassword(account.Password, req.Password) {
		return nil, errs.ErrWrongPassword
	}

	jti := uuid.New()

	tokens, err := s.jwtAuthorizer.Generate(auth.WithSubject(account.ID.String()), auth.WithID(jti.String()))
	if err != nil {
		return nil, err
	}

	t := &entities.Token{
		ID:           jti,
		AccountID:    account.ID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}

	if err = s.tokenRepository.CreateNoReturn(ctx, t); err != nil {
		return nil, err
	}

	session := &entities.Session{
		ID:        jti,
		ExpiresAt: t.ExpiresAt,
	}

	if err = s.sessionService.Set(ctx, session); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthService) Refresh(ctx context.Context, rt string) (*auth.JWT, error) {
	token, ok, err := s.tokenRepository.FindBy(ctx, map[string]interface{}{
		"refresh_token": rt,
	}, nil)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errs.ErrTokenNotFound
	}

	jti := uuid.New()

	tokens, err := s.jwtAuthorizer.Refresh(auth.WithRefreshToken(rt), auth.WithTokenID(jti.String()))
	if err != nil {
		return nil, err
	}

	if err = s.tokenRepository.DeleteAll(ctx, map[string]interface{}{
		"access_token": token.AccessToken,
	}); err != nil {
		return nil, err
	}

	t := &entities.Token{
		ID:           jti,
		AccountID:    token.AccountID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}

	if err = s.tokenRepository.CreateNoReturn(ctx, t); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthService) Logout(ctx context.Context, at string) (err error) {
	token, ok, err := s.tokenRepository.FindBy(ctx, map[string]interface{}{
		"access_token": at,
	}, nil)
	if err != nil {
		return err
	}

	if !ok {
		return errs.ErrTokenNotFound
	}

	tokens, err := s.tokenRepository.Find(ctx, map[string]interface{}{
		"account_id": token.AccountID,
	})

	if err != nil {
		return err
	}

	tokenIDs := lo.Map(tokens, func(item *entities.Token, index int) uuid.UUID {
		return item.ID
	})

	if err := s.tokenRepository.Delete(ctx, tokenIDs, map[string]interface{}{}); err != nil {
		return err
	}

	if err := s.sessionService.Delete(ctx, tokenIDs); err != nil {
		return err
	}

	return nil
}
