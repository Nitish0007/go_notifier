package tests

import (
	"context"
	"testing"

	"github.com/Nitish0007/go_notifier/internal/common/helpers"
	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/apiKey"
	"github.com/stretchr/testify/require"
)

func TestApiKeyRepository_Create_Success(t *testing.T) {
	db, err := helpers.SetupUnitTestsDB()
	require.NoError(t, err, "Failed to setup test database")

	err = helpers.AutoMigrate(db)
	require.NoError(t, err, "Failed to migrate tables")

	repo := apiKey.NewApiKeyRepository(db)
	accountRepo := account.NewAccountRepository(db, repo)

	testAccount := &account.Account{
		Email:             "test@example.com",
		EncryptedPassword: "$2a$10$hashedpassword",
		FirstName:         "Test",
		LastName:          "User",
		IsActive:          true,
	}

	err = accountRepo.Create(context.Background(), testAccount)
	require.NoError(t, err, "Create should succeed")

	testApiKey := &apiKey.ApiKey{
		Key:       "test-key",
		AccountID: testAccount.ID,
	}

	err = repo.Create(context.Background(), testApiKey)
	require.NoError(t, err, "Create should succeed")
	require.Greater(t, testApiKey.ID, 0, "API key should have an ID after creation")
}

func TestApiKeyRepository_Create_DuplicateKeyFails(t *testing.T) {
	db, err := helpers.SetupUnitTestsDB()
	require.NoError(t, err, "Failed to setup test database")

	err = helpers.AutoMigrate(db)
	require.NoError(t, err, "Failed to migrate tables")

	repo := apiKey.NewApiKeyRepository(db)
	accountRepo := account.NewAccountRepository(db, repo)

	testAccount := &account.Account{
		Email:             "dupkey@example.com",
		EncryptedPassword: "$2a$10$hashedpassword",
		FirstName:         "Test",
		LastName:          "User",
		IsActive:          true,
	}
	require.NoError(t, accountRepo.Create(context.Background(), testAccount))

	key := "same-unique-key"
	require.NoError(t, repo.Create(context.Background(), &apiKey.ApiKey{Key: key, AccountID: testAccount.ID}))

	err = repo.Create(context.Background(), &apiKey.ApiKey{Key: key, AccountID: testAccount.ID})
	require.Error(t, err, "Create should fail when key violates unique constraint")
}

func TestApiKeyRepository_FindByAccountID_Success(t *testing.T) {
	db, err := helpers.SetupUnitTestsDB()
	require.NoError(t, err, "Failed to setup test database")

	err = helpers.AutoMigrate(db)
	require.NoError(t, err, "Failed to migrate tables")

	repo := apiKey.NewApiKeyRepository(db)
	accountRepo := account.NewAccountRepository(db, repo)

	testAccount := &account.Account{
		Email:             "test@example.com",
		EncryptedPassword: "$2a$10$hashedpassword",
		FirstName:         "Test",
		LastName:          "User",
		IsActive:          true,
	}

	err = accountRepo.Create(context.Background(), testAccount)
	require.NoError(t, err, "Create should succeed")

	testApiKey := &apiKey.ApiKey{
		Key:       "test-key",
		AccountID: testAccount.ID,
	}

	err = repo.Create(context.Background(), testApiKey)
	require.NoError(t, err, "Create should succeed")

	foundApiKey, err := repo.FindByAccountID(context.Background(), testAccount.ID)
	require.NoError(t, err, "Find by account ID should succeed")
	require.Equal(t, testApiKey.Key, foundApiKey.Key, "API key should be found")
}

func TestApiKeyRepository_FindByAccountID_Failure(t *testing.T) {
	db, err := helpers.SetupUnitTestsDB()
	require.NoError(t, err, "Failed to setup test database")

	err = helpers.AutoMigrate(db)
	require.NoError(t, err, "Failed to migrate tables")

	repo := apiKey.NewApiKeyRepository(db)

	foundApiKey, err := repo.FindByAccountID(context.Background(), 1)
	require.Error(t, err, "Find by account ID should fail with non-existent account ID")
	require.Empty(t, foundApiKey, "API key should be empty")
	require.Equal(t, apiKey.ApiKey{}, foundApiKey, "API key should be empty")
}
