package account

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func InitAccount(payload *SignupRequest) (*Account, error) {
	password := payload.Account.Password
	confirmPassword := payload.Account.ConfirmPassword
	if password != confirmPassword {
		return nil, errors.New("password and confirm password do not match")
	}

	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to encrypt password")
	}

	account := &Account{
		Email: payload.Account.Email,
		EncryptedPassword: string(encryptedPassword),
		FirstName: payload.Account.FirstName,
		LastName: payload.Account.LastName,
		IsActive: true,
	}

	return account, nil
}