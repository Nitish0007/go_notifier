package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Nitish0007/go_notifier/internal/common/helpers"
	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/apiKey"
	"github.com/Nitish0007/go_notifier/internal/features/contact"
	"github.com/Nitish0007/go_notifier/internal/features/emailcontact"
)

func seedAccount(t *testing.T, db *gorm.DB) *account.Account {
	t.Helper()
	apiKeyRepo := apiKey.NewApiKeyRepository(db)
	repo := account.NewAccountRepository(db, apiKeyRepo)
	acc := &account.Account{
		Email:             uuid.NewString() + "@example.com",
		EncryptedPassword: "$2a$10$hashedpassword",
		FirstName:         "A",
		LastName:          "B",
		IsActive:          true,
	}
	require.NoError(t, repo.Create(context.Background(), acc))
	return acc
}

func TestContactRepository_CreateWithEmail_Success(t *testing.T) {
	db, err := helpers.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, helpers.AutoMigrate(db))

	acc := seedAccount(t, db)
	ecRepo := emailcontact.NewEmailContactRepository(db)
	repo := contact.NewContactRepository(db, ecRepo)

	c := &contact.Contact{
		UUID:      uuid.NewString(),
		AccountID: acc.ID,
		FirstName: "Ada",
		LastName:  "Lovelace",
	}
	ec := &emailcontact.EmailContact{
		Email:     "ada@example.com",
		AccountID: acc.ID,
	}

	require.NoError(t, repo.CreateWithEmail(context.Background(), c, ec))
	require.NotZero(t, c.ID)
	require.NotZero(t, ec.ContactID)
	require.Equal(t, c.ID, ec.ContactID)

	loaded, err := repo.FindById(context.Background(), c.ID)
	require.NoError(t, err)
	require.Equal(t, "Ada", loaded.FirstName)
	require.NotNil(t, loaded.EmailContact)
	require.Equal(t, "ada@example.com", loaded.EmailContact.Email)
}

func TestContactRepository_CreateWithEmail_DuplicateEmailPerAccount(t *testing.T) {
	db, err := helpers.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, helpers.AutoMigrate(db))

	acc := seedAccount(t, db)
	ecRepo := emailcontact.NewEmailContactRepository(db)
	repo := contact.NewContactRepository(db, ecRepo)

	email := "same@example.com"
	c1 := &contact.Contact{UUID: uuid.NewString(), AccountID: acc.ID, FirstName: "A", LastName: "1"}
	require.NoError(t, repo.CreateWithEmail(context.Background(), c1, &emailcontact.EmailContact{Email: email, AccountID: acc.ID}))

	c2 := &contact.Contact{UUID: uuid.NewString(), AccountID: acc.ID, FirstName: "B", LastName: "2"}
	err = repo.CreateWithEmail(context.Background(), c2, &emailcontact.EmailContact{Email: email, AccountID: acc.ID})
	require.Error(t, err)
}

func TestContactRepository_GetContacts(t *testing.T) {
	db, err := helpers.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, helpers.AutoMigrate(db))

	acc := seedAccount(t, db)
	ecRepo := emailcontact.NewEmailContactRepository(db)
	repo := contact.NewContactRepository(db, ecRepo)

	for i := 0; i < 2; i++ {
		c := &contact.Contact{UUID: uuid.NewString(), AccountID: acc.ID, FirstName: "F", LastName: string(rune('A' + i))}
		require.NoError(t, repo.CreateWithEmail(context.Background(), c, &emailcontact.EmailContact{
			Email:     uuid.NewString() + "@example.com",
			AccountID: acc.ID,
		}))
	}

	list, err := repo.GetContacts(context.Background(), acc.ID)
	require.NoError(t, err)
	require.Len(t, list, 2)
	for _, row := range list {
		require.NotNil(t, row.EmailContact)
	}
}

func TestContactRepository_FindById_NotFound(t *testing.T) {
	db, err := helpers.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, helpers.AutoMigrate(db))

	repo := contact.NewContactRepository(db, emailcontact.NewEmailContactRepository(db))
	_, err = repo.FindById(context.Background(), 99999)
	require.Error(t, err)
}

func TestContactRepository_FindByUUID(t *testing.T) {
	db, err := helpers.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, helpers.AutoMigrate(db))

	acc := seedAccount(t, db)
	ecRepo := emailcontact.NewEmailContactRepository(db)
	repo := contact.NewContactRepository(db, ecRepo)

	u := uuid.NewString()
	c := &contact.Contact{UUID: u, AccountID: acc.ID, FirstName: "U", LastName: "UID"}
	require.NoError(t, repo.CreateWithEmail(context.Background(), c, &emailcontact.EmailContact{
		Email:     "uid@example.com",
		AccountID: acc.ID,
	}))

	found, err := repo.FindByUUID(context.Background(), u)
	require.NoError(t, err)
	require.Equal(t, c.ID, found.ID)
	require.NotNil(t, found.EmailContact)
}
