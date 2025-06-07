package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Nitish0007/go_notifier/internal/models"
	"github.com/Nitish0007/go_notifier/internal/repositories"
	"github.com/Nitish0007/go_notifier/utils"
)

type AccountService struct {
	AccRepo   *repositories.AccountRepository
	ApiKeyRepo *repositories.ApiKeyRepository
}

func NewAccountService(accRepo *repositories.AccountRepository, apiKeyRepo *repositories.ApiKeyRepository) *AccountService {
	return &AccountService{
		AccRepo: accRepo,
		ApiKeyRepo: apiKeyRepo,
	}
}

// PUBLIC METHODS with receiver
func (s *AccountService) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	if err := validateAccount(account); err != nil {
		return nil, err
	}

	// Start a transaction
	tx, err := s.AccRepo.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// create account
	if err := s.AccRepo.CreateTx(ctx, account, tx); err != nil {
		return nil, err
	}

	// initialize API key
	apiKey := &models.ApiKey{
		Key:      utils.GenerateAlphaNumericKey(),
		AccountID: account.ID,
	}

	// create API key
	if err := s.ApiKeyRepo.CreateTx(ctx, apiKey, tx); err != nil {
		return nil, err
	}

	// NOTE: TODO for multi-tenancy support
	// Can be created a new DB tenant with name 'account_${account.ID}' if one account has huge amount of data
	// Tenancy scheme can be implemented here: single database schema-based multi-tenancy
	// migration needs to be handled separately for each tenant

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) InitializeAccount(ctx context.Context, accountData map[string]any) (*models.Account, error) {
	account := &models.Account{}

	if accountData != nil {
		if firstName, ok := accountData["first_name"].(string); ok {
			account.FirstName = firstName
		}

		if lastName, ok := accountData["last_name"].(string); ok {
			account.LastName = lastName
		}

		if email, ok := accountData["email"].(string); ok {
			email = strings.ToLower(email)
			if !utils.ValidateEmail(email) {
				return nil, errors.New("invalid email format")
			}
			account.Email = email
		}

		password, ok := accountData["password"].(string)
		confirmPass, ok2 := accountData["confirm_password"].(string)

		if !ok || !ok2 || password == "" || confirmPass == "" {
			return nil, errors.New("password and confirm password")
		}
		if password != confirmPass {
			return nil, errors.New("password and confirm password do not match")
		}

		if len(password) < 6 {
			return nil, errors.New("password must be at least 6 characters long")
		}
		encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to encrypt password")
		}
		account.EncryptedPassword = string(encryptedPass)

		isActive, ok := accountData["is_active"].(bool)
		if !ok {
			account.IsActive = true // default to true if not provided
		} else {
			account.IsActive = isActive
		}

		account.CreatedAt = time.Now()
		account.UpdatedAt = time.Now()

		return account, nil
	}

	return nil, errors.New("account data not provided")
}

func (s * AccountService) Login(ctx context.Context, payload map[string]any) (string, error) {
	// validating payload
	email, exists := payload["email"].(string)
	if !exists || email == "" {
		return "", errors.New("email not provided")
	}
	email = strings.ToLower(email)

	if !utils.ValidateEmail(email) {
		return "", errors.New("invalid email format")
	}

	// find by email
	account, err := s.AccRepo.FindAccountByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(account.EncryptedPassword), []byte(payload["password"].(string)))
	if err != nil {
		return "", err
	}
	
	apiKey, err := s.ApiKeyRepo.FindByAccountID(ctx, account.ID)
	
	if err != nil {
		return "", err
	}

	return apiKey.Key, nil
}

// PRIVATE METHODS without receiver 
func validateAccount(account *models.Account) error {
	if account.FirstName == "" || account.Email == "" || account.EncryptedPassword == "" {
		return errors.New("invalid account parameters: first name, email, and encrypted password are required")
	}
	return nil
}