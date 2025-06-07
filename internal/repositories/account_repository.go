package repositories

import (
	"context"
	"github.com/Nitish0007/go_notifier/internal/models"

	"github.com/jackc/pgx/v5"
)

type AccountRepository struct {
	DB *pgx.Conn
}

func NewAccountRepository(conn *pgx.Conn) *AccountRepository {
	return &AccountRepository{
		DB: conn,
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

func (r *AccountRepository) FindAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	query := `SELECT * FROM accounts WHERE email LIKE $1`
	account := models.Account{}
	err := r.DB.QueryRow(ctx, query, email).Scan(
		&account.ID,
		&account.Email,
		&account.EncryptedPassword,
		&account.FirstName,
		&account.LastName,
		&account.IsActive,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if(err != nil){
		return nil, err
	}
	return &account, nil
}