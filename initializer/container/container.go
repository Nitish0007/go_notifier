package container

import (
	"gorm.io/gorm"

	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/apiKey"
	"github.com/Nitish0007/go_notifier/internal/features/configuration"
	"github.com/Nitish0007/go_notifier/internal/features/contact"
	"github.com/Nitish0007/go_notifier/internal/features/emailcontact"
	"github.com/Nitish0007/go_notifier/internal/features/emailnotification"

	// "github.com/Nitish0007/go_notifier/internal/lib/notifier"
	// notifierInterface "github.com/Nitish0007/go_notifier/internal/shared/interfaces/notifier"
)

type Container struct {
	DB *gorm.DB
	// NOTE: keep in mind that dependencies should be intialized first
	// repositories are the lowest level of dependencies and should be intialized first

	// Repositories
	ApiKeyRepository        *apiKey.ApiKeyRepository
	AccountRepository       *account.AccountRepository
	ConfigurationRepository *configuration.ConfigurationRepository
	EmailNotificationRepository  *emailnotification.EmailNotificationRepository
	ContactRepository       *contact.ContactRepository
	EmailContactRepository  *emailcontact.EmailContactRepository

	// Notifiers
	// Notifiers []notifierInterface.Notifier

	// Services
	AccountService       *account.AccountService
	ConfigurationService *configuration.ConfigurationService
	EmailNotificationService  *emailnotification.EmailNotificationService
	ContactService       *contact.ContactService

	// Handlers
	AccountHandler       *account.AccountHandler
	ConfigurationHandler *configuration.ConfigurationHandler
	EmailNotificationHandler  *emailnotification.EmailNotificationHandler
	ContactHandler       *contact.ContactHandler
}

func NewContainer(db *gorm.DB) *Container {
	c := &Container{DB: db}

	c.InitializeRepositories()
	// c.InitializeNotifiers()
	c.InitializeServices()
	c.InitializeHandlers()
	return c
}

// func (c *Container) InitializeNotifiers() {
// 	c.Notifiers = []notifierInterface.Notifier{
// 		notifier.NewEmailNotifier(c.EmailNotificationRepository),
// 	}
// }

func (c *Container) InitializeRepositories() {
	// intialize shared repositories first
	c.ApiKeyRepository = apiKey.NewApiKeyRepository(c.DB)

	c.AccountRepository = account.NewAccountRepository(c.DB, c.ApiKeyRepository)
	c.ConfigurationRepository = configuration.NewConfigurationRepository(c.DB)
	c.ContactRepository = contact.NewContactRepository(c.DB, c.EmailContactRepository)
	c.EmailContactRepository = emailcontact.NewEmailContactRepository(c.DB)
	c.EmailNotificationRepository = emailnotification.NewEmailNotificationRepository(c.DB)
}

func (c *Container) InitializeServices() {
	c.AccountService = account.NewAccountService(c.AccountRepository)
	c.ConfigurationService = configuration.NewConfigurationService(c.ConfigurationRepository)
	c.ContactService = contact.NewContactService(c.ContactRepository)
	c.EmailNotificationService = emailnotification.NewEmailNotificationService(c.EmailNotificationRepository)
}

func (c *Container) InitializeHandlers() {
	c.AccountHandler = account.NewAccountHandler(c.AccountService)
	c.ConfigurationHandler = configuration.NewConfigurationHandler(c.ConfigurationService)
	c.ContactHandler = contact.NewContactHandler(c.ContactService)
	c.EmailNotificationHandler = emailnotification.NewEmailNotificationHandler(c.EmailNotificationService)
}
