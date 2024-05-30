package entities

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrPasswordsMismatch = errors.New("password mismatch")

var RootLogin = "root"

type Account struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Login    string    `json:"login"`
	Password string    `json:"-"`
}

func NewAccount(req CreateAccount) (*Account, error) {
	if req.Password != req.PasswordConfirm {
		return nil, ErrPasswordsMismatch
	}

	passwordHash, err := GeneratePasswordHash(req.Password)
	if err != nil {
		return nil, err
	}

	a := &Account{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      req.Name,
		Login:     req.Login,
		Password:  passwordHash,
	}

	return a, nil
}

func (a *Account) TableName() string {
	return "accounts"
}

func (a *Account) GetID() uuid.UUID {
	return a.ID
}

func (a *Account) IDColumnName() string {
	return "id"
}

type DeleteAccountRequest struct {
	ID uuid.UUID `json:"id" validate:"required,not_root_id"`
}

type CreateAccount struct {
	Name            string `json:"name" validate:"required"`
	Login           string `json:"login" validate:"required,unique_login"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}

const bcryptCost = 15

func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPassword(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}

	return true
}
