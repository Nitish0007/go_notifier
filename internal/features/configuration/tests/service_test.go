package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Nitish0007/go_notifier/internal/common/database"
	"github.com/Nitish0007/go_notifier/internal/features/configuration"
)

func boolPtr(b bool) *bool { return &b }

func TestConfigurationService_CreateAndList(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	acc := seedAccount(t, db)
	svc := configuration.NewConfigurationService(configuration.NewConfigurationRepository(db))

	req := &configuration.CreateConfigurationRequest{}
	req.Configuration.AccountID = acc.ID
	req.Configuration.DefaultConfiguration = boolPtr(false)
	req.Configuration.ConfigType = "smtp"
	req.Configuration.ConfigurationData = map[string]any{
		"host": "smtp.example.com", "port": float64(587),
		"username": "u", "password": "p", "from": "from@example.com",
	}

	created, err := svc.CreateConfiguration(context.Background(), req)
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	list, err := svc.GetConfigurations(context.Background(), acc.ID)
	require.NoError(t, err)
	require.Len(t, list, 1)
}

func TestConfigurationService_UpdateConfiguration(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	acc := seedAccount(t, db)
	svc := configuration.NewConfigurationService(configuration.NewConfigurationRepository(db))

	createReq := &configuration.CreateConfigurationRequest{}
	createReq.Configuration.AccountID = acc.ID
	createReq.Configuration.DefaultConfiguration = boolPtr(false)
	createReq.Configuration.ConfigType = "smtp"
	createReq.Configuration.ConfigurationData = map[string]any{
		"host": "old.example.com", "port": float64(25),
		"username": "u", "password": "p", "from": "from@example.com",
	}

	created, err := svc.CreateConfiguration(context.Background(), createReq)
	require.NoError(t, err)

	upd := &configuration.UpdateConfigurationRequest{}
	upd.Configuration.ID = created.ID
	upd.Configuration.AccountID = acc.ID
	upd.Configuration.ConfigType = "smtp"
	upd.Configuration.DefaultConfiguration = boolPtr(false)
	upd.Configuration.ConfigurationData = map[string]any{
		"host": "new.example.com", "port": float64(465),
		"username": "u2", "password": "p2", "from": "from2@example.com",
	}

	updated, err := svc.UpdateConfiguration(context.Background(), upd)
	require.NoError(t, err)
	require.Equal(t, created.ID, updated.ID)
	require.Equal(t, "new.example.com", updated.ConfigurationData["host"])
}

func TestConfigurationService_DeleteConfiguration(t *testing.T) {
	db, err := database.SetupUnitTestsDB()
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(db))

	acc := seedAccount(t, db)
	repo := configuration.NewConfigurationRepository(db)
	svc := configuration.NewConfigurationService(repo)

	createReq := &configuration.CreateConfigurationRequest{}
	createReq.Configuration.AccountID = acc.ID
	createReq.Configuration.DefaultConfiguration = boolPtr(false)
	createReq.Configuration.ConfigType = "in_app"
	createReq.Configuration.ConfigurationData = map[string]any{
		"web_app_url": "https://app.example", "web_app_secret": "s", "web_app_token": "t",
	}

	created, err := svc.CreateConfiguration(context.Background(), createReq)
	require.NoError(t, err)

	require.NoError(t, svc.DeleteConfiguration(context.Background(), acc.ID, created.ID))

	_, err = repo.GetByFields(context.Background(), map[string]any{"id": created.ID})
	require.Error(t, err)
}
