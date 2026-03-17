package helpers

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/api_key"
	"github.com/Nitish0007/go_notifier/utils"
)

// SetupTestDBForIntegrationTests connects to test database and runs migrations
// Returns the GORM DB instance for integration tests
func SetupTestDBForIntegrationTests() (*gorm.DB, error) {
	return utils.SetupTestDB()
}

// SetupUnitTestsDB creates an in-memory SQLite database for unit tests
// Returns the GORM DB instance
func SetupUnitTestsDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// AutoMigrate runs GORM AutoMigrate on the provided database
// Migrates Account and ApiKey models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&account.Account{},
		&api_key.ApiKey{},
	)
}
