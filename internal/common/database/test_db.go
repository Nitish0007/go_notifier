package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/apiKey"
	"github.com/Nitish0007/go_notifier/internal/features/configuration"
	"github.com/Nitish0007/go_notifier/internal/features/contact"
	"github.com/Nitish0007/go_notifier/internal/features/emailcontact"
)

// SetupTestDBForIntegrationTests connects to test database and runs migrations
// Returns the GORM DB instance for integration tests
func SetupTestDBForIntegrationTests() (*gorm.DB, error) {
	return SetupTestDB()
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

// AutoMigrate runs GORM AutoMigrate on the provided database for in-memory SQLite tests.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&account.Account{},
		&apiKey.ApiKey{},
		&contact.Contact{},
		&emailcontact.EmailContact{},
		&configuration.Configuration{},
	)
}
