package repositories

import (
	"context"
	"github.com/Nitish0007/go_notifier/internal/models"
	"github.com/Nitish0007/go_notifier/utils"

	"github.com/jackc/pgx/v5"
)

type AccountRepository struct {
	DB *pgx.Conn
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		DB: utils.DB,
	}
}

func (r *AccountRepository) Create(ctx context.Context, account *models.Account) error {
	query := `INSERT INTO accounts (first_name, last_name, email, encrypted_password) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.DB.QueryRow(ctx, query, account.FirstName, account.LastName, account.Email, account.EncryptedPassword).Scan(&account.ID)
	if err != nil {
		return err
	}
	return nil
}

// CreateTx creates a new account within a transaction.
func (r *AccountRepository) CreateTx(ctx context.Context, account *models.Account, tx pgx.Tx) error {
	query := `INSERT INTO accounts (first_name, last_name, email, encrypted_password) VALUES ($1, $2, $3, $4) RETURNING id`
	err := tx.QueryRow(ctx, query, account.FirstName, account.LastName, account.Email, account.EncryptedPassword).Scan(&account.ID)
	if err != nil {
		return err
	}
	return nil
}