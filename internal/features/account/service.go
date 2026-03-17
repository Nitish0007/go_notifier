package account

import (
	"strings"
	"context"
	"strconv"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	accountRepository *AccountRepository
}

func NewAccountService(accountRepository *AccountRepository) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, payload *SignupRequest) (*SignupResponse, error) {
	account, err := InitAccount(payload)
	if err != nil {
		return nil, err
	}

	if err := s.accountRepository.RegisterAccount(ctx, account); err != nil {
		return nil, err
	}

	return &SignupResponse{
		ID: account.ID,
		Email: account.Email,
		FirstName: account.FirstName,
		LastName: account.LastName,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}, nil
}

func (s *AccountService) Login(ctx context.Context, payload *LoginRequest) (*LoginResponse, error) {
	email := payload.Login.Email
	password := payload.Login.Password
	email = strings.ToLower(email)

	account, err := s.accountRepository.FindAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.EncryptedPassword), []byte(password))
	if err != nil {
		return nil, err
	}

	apiKey, err := s.accountRepository.GetApiKeyByAccountID(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AuthToken: apiKey,
		AccountID: strconv.Itoa(account.ID),
		Email: account.Email,
		FirstName: account.FirstName,
		LastName: account.LastName,
	}, nil
}