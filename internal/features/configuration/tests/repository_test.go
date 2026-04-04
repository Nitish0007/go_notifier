package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Nitish0007/go_notifier/internal/common/database"
	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/apiKey"
	"github.com/Nitish0007/go_notifier/internal/features/configuration"
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

func TestConfigurationRepository_Create_And_Index(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	acc := seedAccount(t, db)
	repo := configuration.NewConfigurationRepository(db)

	cfg := &configuration.Configuration{
		AccountID:            acc.ID,
		DefaultConfiguration: false,
		ConfigType:           "smtp",
		ConfigurationData: map[string]any{
			"host": "smtp.example.com", "port": float64(587),
			"username": "u", "password": "p", "from": "from@example.com",
		},
	}
	require.NoError(t, repo.Create(context.Background(), cfg))
	require.NotZero(t, cfg.ID)

	list, err := repo.Index(context.Background(), acc.ID)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "smtp", list[0].ConfigType)
}

func TestConfigurationRepository_GetByAccountID_DefaultTrue(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	acc := seedAccount(t, db)
	repo := configuration.NewConfigurationRepository(db)

	row := &configuration.Configuration{
		AccountID:            acc.ID,
		DefaultConfiguration: true,
		ConfigType:           "smtp",
		ConfigurationData:    map[string]any{"host": "h", "port": float64(1)},
	}
	require.NoError(t, db.Session(&gorm.Session{SkipHooks: true}).Create(row).Error)

	found, err := repo.GetByAccountID(context.Background(), acc.ID)
	require.NoError(t, err)
	require.True(t, found.DefaultConfiguration)
	require.Equal(t, acc.ID, found.AccountID)
}

func TestConfigurationRepository_GetByAccountIDAndConfigType(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	acc := seedAccount(t, db)
	repo := configuration.NewConfigurationRepository(db)

	require.NoError(t, repo.Create(context.Background(), &configuration.Configuration{
		AccountID:            acc.ID,
		DefaultConfiguration: false,
		ConfigType:           "smtp",
		ConfigurationData:    map[string]any{"host": "x", "port": float64(25)},
	}))

	found, err := repo.GetByAccountIDAndConfigType(context.Background(), acc.ID, "smtp")
	require.NoError(t, err)
	require.Equal(t, "smtp", found.ConfigType)
	require.False(t, found.DefaultConfiguration)
}

func TestConfigurationRepository_Update(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	acc := seedAccount(t, db)
	repo := configuration.NewConfigurationRepository(db)

	cfg := &configuration.Configuration{
		AccountID:            acc.ID,
		DefaultConfiguration: false,
		ConfigType:           "smtp",
		ConfigurationData:    map[string]any{"host": "old", "port": float64(25)},
	}
	require.NoError(t, repo.Create(context.Background(), cfg))

	cfg.ConfigurationData = map[string]any{"host": "newhost", "port": float64(587)}
	require.NoError(t, repo.Update(context.Background(), cfg))

	reloaded, err := repo.GetByFields(context.Background(), map[string]any{"id": cfg.ID})
	require.NoError(t, err)
	require.Equal(t, "newhost", reloaded.ConfigurationData["host"])
}

func TestConfigurationRepository_Delete(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	acc := seedAccount(t, db)
	repo := configuration.NewConfigurationRepository(db)

	cfg := &configuration.Configuration{
		AccountID:            acc.ID,
		DefaultConfiguration: false,
		ConfigType:           "in_app",
		ConfigurationData:    map[string]any{"web_app_url": "https://a.example", "web_app_secret": "s", "web_app_token": "t"},
	}
	require.NoError(t, repo.Create(context.Background(), cfg))

	require.NoError(t, repo.Delete(context.Background(), cfg.ID))
	err = repo.Delete(context.Background(), cfg.ID)
	require.Error(t, err)
}

func TestConfigurationRepository_GetByFields(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	acc := seedAccount(t, db)
	repo := configuration.NewConfigurationRepository(db)

	cfg := &configuration.Configuration{
		AccountID:            acc.ID,
		DefaultConfiguration: false,
		ConfigType:           "smtp",
		ConfigurationData:    map[string]any{"host": "byfields"},
	}
	require.NoError(t, repo.Create(context.Background(), cfg))

	found, err := repo.GetByFields(context.Background(), map[string]any{
		"account_id":  acc.ID,
		"config_type": "smtp",
	})
	require.NoError(t, err)
	require.Equal(t, cfg.ID, found.ID)
}
