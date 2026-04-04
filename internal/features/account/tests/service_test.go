package tests

import (
	"context"
	"testing"

	"github.com/Nitish0007/go_notifier/internal/common/database"
	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/apiKey"
	"github.com/stretchr/testify/require"
)

func TestAccountService_CreateAccount_Success(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	apiKeyRepo := apiKey.NewApiKeyRepository(db)
	accRepo := account.NewAccountRepository(db, apiKeyRepo)
	svc := account.NewAccountService(accRepo)

	req := &account.SignupRequest{}
	req.Account.Email = "signup@example.com"
	req.Account.Password = "secret123"
	req.Account.ConfirmPassword = "secret123"
	req.Account.FirstName = "Sam"
	req.Account.LastName = "Sign"

	resp, err := svc.CreateAccount(context.Background(), req)
	require.NoError(t, err)
	require.NotZero(t, resp.ID)
	require.Equal(t, "signup@example.com", resp.Email)
	require.Equal(t, "Sam", resp.FirstName)

	_, err = apiKeyRepo.FindByAccountID(context.Background(), resp.ID)
	require.NoError(t, err)
}

func TestAccountService_CreateAccount_PasswordMismatch(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	svc := account.NewAccountService(account.NewAccountRepository(db, apiKey.NewApiKeyRepository(db)))

	req := &account.SignupRequest{}
	req.Account.Email = "x@example.com"
	req.Account.Password = "secret123"
	req.Account.ConfirmPassword = "other"
	req.Account.FirstName = "A"
	req.Account.LastName = "B"

	_, err = svc.CreateAccount(context.Background(), req)
	require.Error(t, err)
}

func TestAccountService_Login_Success(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	apiKeyRepo := apiKey.NewApiKeyRepository(db)
	accRepo := account.NewAccountRepository(db, apiKeyRepo)
	svc := account.NewAccountService(accRepo)

	signup := &account.SignupRequest{}
	signup.Account.Email = "login@example.com"
	signup.Account.Password = "mypassword123"
	signup.Account.ConfirmPassword = "mypassword123"
	signup.Account.FirstName = "L"
	signup.Account.LastName = "User"
	_, err = svc.CreateAccount(context.Background(), signup)
	require.NoError(t, err)

	login := &account.LoginRequest{}
	login.Login.Email = "LOGIN@EXAMPLE.COM"
	login.Login.Password = "mypassword123"

	out, err := svc.Login(context.Background(), login)
	require.NoError(t, err)
	require.NotEmpty(t, out.AuthToken)
	require.Equal(t, "login@example.com", out.Email)
}

func TestAccountService_Login_WrongPassword(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	apiKeyRepo := apiKey.NewApiKeyRepository(db)
	accRepo := account.NewAccountRepository(db, apiKeyRepo)
	svc := account.NewAccountService(accRepo)

	signup := &account.SignupRequest{}
	signup.Account.Email = "u@example.com"
	signup.Account.Password = "correcthorse"
	signup.Account.ConfirmPassword = "correcthorse"
	signup.Account.FirstName = "U"
	signup.Account.LastName = "V"
	_, err = svc.CreateAccount(context.Background(), signup)
	require.NoError(t, err)

	login := &account.LoginRequest{}
	login.Login.Email = "u@example.com"
	login.Login.Password = "wrong"

	_, err = svc.Login(context.Background(), login)
	require.Error(t, err)
}
