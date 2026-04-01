package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Nitish0007/go_notifier/internal/common/helpers"
	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/apiKey"
	"github.com/Nitish0007/go_notifier/internal/features/contact"
	"github.com/Nitish0007/go_notifier/internal/features/emailcontact"
)

func newContactService(t *testing.T) (*account.Account, *contact.ContactService) {
	t.Helper()
	db, err := helpers.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, helpers.AutoMigrate(db))

	apiKeyRepo := apiKey.NewApiKeyRepository(db)
	accRepo := account.NewAccountRepository(db, apiKeyRepo)
	acc := &account.Account{
		Email:             uuid.NewString() + "@example.com",
		EncryptedPassword: "$2a$10$hashedpassword",
		FirstName:         "A",
		LastName:          "B",
		IsActive:          true,
	}
	require.NoError(t, accRepo.Create(context.Background(), acc))

	ecRepo := emailcontact.NewEmailContactRepository(db)
	cr := contact.NewContactRepository(db, ecRepo)
	return acc, contact.NewContactService(cr)
}

func TestContactService_CreateContact_Success(t *testing.T) {
	acc, svc := newContactService(t)

	req := &contact.CreateContactRequest{}
	req.Contact.AccountID = acc.ID
	req.Contact.FirstName = "Pat"
	req.Contact.LastName = "Lee"
	req.Contact.Email = "pat@example.com"

	resp, err := svc.CreateContact(context.Background(), req)
	require.NoError(t, err)
	require.NotZero(t, resp.ID)
	require.NotEmpty(t, resp.UUID)
	require.Equal(t, "pat@example.com", resp.Email)
	require.Equal(t, acc.ID, resp.AccountID)
}

func TestContactService_GetContacts(t *testing.T) {
	acc, svc := newContactService(t)

	for i := 0; i < 2; i++ {
		req := &contact.CreateContactRequest{}
		req.Contact.AccountID = acc.ID
		req.Contact.FirstName = "F"
		req.Contact.LastName = string(rune('X' + i))
		req.Contact.Email = uuid.NewString() + "@example.com"
		_, err := svc.CreateContact(context.Background(), req)
		require.NoError(t, err)
	}

	list, err := svc.GetContacts(context.Background(), acc.ID)
	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestContactService_GetContactByKey_ID(t *testing.T) {
	acc, svc := newContactService(t)

	req := &contact.CreateContactRequest{}
	req.Contact.AccountID = acc.ID
	req.Contact.FirstName = "Q"
	req.Contact.LastName = "R"
	req.Contact.Email = "qr@example.com"
	created, err := svc.CreateContact(context.Background(), req)
	require.NoError(t, err)

	out, err := svc.GetContactByKey(context.Background(), "id", created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, out.ID)
	require.Equal(t, "qr@example.com", out.Email)
}

func TestContactService_GetContactByKey_UUID(t *testing.T) {
	acc, svc := newContactService(t)

	req := &contact.CreateContactRequest{}
	req.Contact.AccountID = acc.ID
	req.Contact.FirstName = "M"
	req.Contact.LastName = "N"
	req.Contact.Email = "mn@example.com"
	created, err := svc.CreateContact(context.Background(), req)
	require.NoError(t, err)
	require.NotEmpty(t, created.UUID)

	out, err := svc.GetContactByKey(context.Background(), "uuid", created.UUID)
	require.NoError(t, err)
	require.Equal(t, created.ID, out.ID)
	require.Equal(t, "mn@example.com", out.Email)
}

func TestContactService_GetContactByKey_InvalidKey(t *testing.T) {
	_, svc := newContactService(t)
	_, err := svc.GetContactByKey(context.Background(), "unknown", 1)
	require.Error(t, err)
}

func TestContactService_CreateContact_UserProvidedUUID(t *testing.T) {
	acc, svc := newContactService(t)
	fixed := uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")

	req := &contact.CreateContactRequest{}
	req.Contact.AccountID = acc.ID
	req.Contact.FirstName = "Custom"
	req.Contact.LastName = "UUID"
	req.Contact.Email = "custom-uuid@example.com"
	req.Contact.UUID = fixed.String()

	resp, err := svc.CreateContact(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, fixed.String(), resp.UUID)
}

func TestContactService_CreateContact_InvalidUUID(t *testing.T) {
	acc, svc := newContactService(t)

	req := &contact.CreateContactRequest{}
	req.Contact.AccountID = acc.ID
	req.Contact.FirstName = "Bad"
	req.Contact.LastName = "UUID"
	req.Contact.Email = "bad-uuid@example.com"
	req.Contact.UUID = "not-a-valid-uuid"

	_, err := svc.CreateContact(context.Background(), req)
	require.Error(t, err)
}
