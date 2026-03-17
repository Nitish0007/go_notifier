package tests

import (
	"context"
	"testing"

	"github.com/Nitish0007/go_notifier/internal/common/helpers"
	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/api_key"
	"github.com/stretchr/testify/require"
)

func TestAccountRepository_Create_Success(t *testing.T) {
	// ARRANGE: Setup test database
	db, err := helpers.SetupUnitTestsDB()
	require.NoError(t, err, "Failed to setup test database")

	// Migrate account and api key tables
	err = helpers.AutoMigrate(db)
	require.NoError(t, err, "Failed to migrate tables")

	// Create repositories
	apiKeyRepo := api_key.NewApiKeyRepository(db)
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

func TestAccountRepository_Create_Failure(t *testing.T) {
	db, err := helpers.SetupUnitTestsDB()
	require.NoError(t, err, "Failed to setup test database")

	err = helpers.AutoMigrate(db)
	require.NoError(t, err, "Faile")
}
