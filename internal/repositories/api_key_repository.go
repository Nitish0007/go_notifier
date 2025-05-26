package repositories

import (
	"context"
	"github.com/jackc/pgx/v5"

	"github.com/Nitish0007/go_notifier/internal/models"
	"github.com/Nitish0007/go_notifier/utils"
)


type ApiKeyRepository struct {
	DB *pgx.Conn
}

func NewApiKeyRepository() *ApiKeyRepository{
	return &ApiKeyRepository{
		DB: utils.DB,
	}
}

func (r *ApiKeyRepository) Create(ctx context.Context, apiKey *models.ApiKey) error {
	query := `INSERT INTO api_keys (key, account_id) VALUES ($1, $2) RETURNING id`
	err := r.DB.QueryRow(ctx, query, apiKey.Key, apiKey.AccountID).Scan(&apiKey.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ApiKeyRepository) CreateTx(ctx context.Context, apiKey *models.ApiKey, tx pgx.Tx) error {
	query := `INSERT INTO api_keys (key, account_id) VALUES ($1, $2) RETURNING id`
	err := tx.QueryRow(ctx, query, apiKey.Key, apiKey.AccountID).Scan(&apiKey.ID)
	if err != nil {
		return err
	}
	return nil
}