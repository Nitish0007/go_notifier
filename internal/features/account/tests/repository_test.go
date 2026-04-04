package tests

import (
	"context"
	"testing"

	"github.com/Nitish0007/go_notifier/internal/common/database"
	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/apiKey"
	"github.com/stretchr/testify/require"
)

func TestAccountRepository_Create_Success(t *testing.T) {
	// ARRANGE: Setup test database
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err, "Failed to setup test database")

	// Migrate account and api key tables
	err = database.AutoMigrate(db)
	require.NoError(t, err, "Failed to migrate tables")

	// Create repositories
	apiKeyRepo := apiKey.NewApiKeyRepository(db)
	repo := account.NewAccountRepository(db, apiKeyRepo)

	// ACT: Create account
	testAccount := &account.Account{
		Email:             "test@example.com",
		EncryptedPassword: "$2a$10$hashedpassword", // Mock encrypted password
		FirstName:         "Test",
		LastName:          "User",
		IsActive:          true,
	}

	err = repo.Create(context.Background(), testAccount)

	// ASSERT
	require.NoError(t, err, "Create should succeed")
	require.Greater(t, testAccount.ID, 0, "Account should have an ID after creation")
}

func TestAccountRepository_Create_DuplicateEmailFails(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err, "Failed to setup test database")

	err = database.AutoMigrate(db)
	require.NoError(t, err, "Failed to migrate tables")

	apiKeyRepo := apiKey.NewApiKeyRepository(db)
	repo := account.NewAccountRepository(db, apiKeyRepo)

	acc := &account.Account{
		Email:             "dup@example.com",
		EncryptedPassword: "$2a$10$hashedpassword",
		FirstName:         "A",
		LastName:          "B",
		IsActive:          true,
	}
	require.NoError(t, repo.Create(context.Background(), acc))

	err = repo.Create(context.Background(), &account.Account{
		Email:             "dup@example.com",
		EncryptedPassword: "$2a$10$other",
		FirstName:         "C",
		LastName:          "D",
		IsActive:          true,
	})
	require.Error(t, err, "Create should fail on duplicate email")
}

func TestAccountRepository_RegisterAccount_Success(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err, "Failed to setup test database")

	err = database.AutoMigrate(db)
	require.NoError(t, err, "Failed to migrate tables")

	apiKeyRepo := apiKey.NewApiKeyRepository(db)
	accRepo := account.NewAccountRepository(db, apiKeyRepo)

	testAccount := &account.Account{
		Email:             "test@example.com",
		EncryptedPassword: "$2a$10$hashedpassword",
		FirstName:         "Test",
		LastName:          "User",
		IsActive:          true,
	}

	err = accRepo.RegisterAccount(context.Background(), testAccount)
	if err != nil {
		t.Fatalf("Register Account should succeed: %v", err)
	}

	apiKey, err := apiKeyRepo.FindByAccountID(context.Background(), testAccount.ID)
	if err != nil {
		t.Fatalf("Find by account ID should succeed: %v", err)
	}
	require.NotEmpty(t, apiKey, "API key should not be empty")
	require.Equal(t, testAccount.ID, apiKey.AccountID, "API key should be found")
}

func TestAccountRepository_FindAccountByEmail_Success(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)

	require.NoError(t, database.AutoMigrate(db))

	apiKeyRepo := apiKey.NewApiKeyRepository(db)
	repo := account.NewAccountRepository(db, apiKeyRepo)

	acc := &account.Account{
		Email:             "find@example.com",
		EncryptedPassword: "$2a$10$hashedpassword",
		FirstName:         "Jane",
		LastName:          "Doe",
		IsActive:          true,
	}
	require.NoError(t, repo.Create(context.Background(), acc))

	found, err := repo.FindAccountByEmail(context.Background(), "find@example.com")
	require.NoError(t, err)
	require.Equal(t, acc.ID, found.ID)
	require.Equal(t, "Jane", found.FirstName)
}

func TestAccountRepository_FindAccountByEmail_NotFound(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	apiKeyRepo := apiKey.NewApiKeyRepository(db)
	repo := account.NewAccountRepository(db, apiKeyRepo)

	_, err = repo.FindAccountByEmail(context.Background(), "missing@example.com")
	require.Error(t, err)
}

func TestAccountRepository_GetApiKeyByAccountID_Success(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	apiKeyRepo := apiKey.NewApiKeyRepository(db)
	repo := account.NewAccountRepository(db, apiKeyRepo)

	acc := &account.Account{
		Email:             "key@example.com",
		EncryptedPassword: "$2a$10$hashedpassword",
		FirstName:         "K",
		LastName:          "L",
		IsActive:          true,
	}
	require.NoError(t, repo.RegisterAccount(context.Background(), acc))

	keyStr, err := repo.GetApiKeyByAccountID(context.Background(), acc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, keyStr)
}
