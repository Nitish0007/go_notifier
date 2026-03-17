package container

import (
	"gorm.io/gorm"

	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/api_key"
	"github.com/Nitish0007/go_notifier/internal/features/configuration"
	"github.com/Nitish0007/go_notifier/internal/features/notification"

	notifierInterface "github.com/Nitish0007/go_notifier/internal/shared/interfaces/notifier"
	"github.com/Nitish0007/go_notifier/internal/lib/notifier"
)

type Container struct {
	DB *gorm.DB
	// NOTE: keep in mind that dependencies should be intialized first
	// repositories are the lowest level of dependencies and should be intialized first

	// Repositories
	ApiKeyRepository 	     			*api_key.ApiKeyRepository
	AccountRepository      			*account.AccountRepository
	ConfigurationRepository 		*configuration.ConfigurationRepository
	NotificationRepository 			*notification.NotificationRepository

	// Notifiers
	Notifiers 							[]notifierInterface.Notifier
	
	// Services
	AccountService      			*account.AccountService
	ConfigurationService 			*configuration.ConfigurationService
	NotificationService 			*notification.NotificationService

	// Handlers
	AccountHandler      			*account.AccountHandler
	ConfigurationHandler 			*configuration.ConfigurationHandler
	NotificationHandler 			*notification.NotificationHandler
}

func NewContainer(db *gorm.DB) *Container {
	c := &Container{DB: db}
	
	c.InitializeRepositories()
	c.InitializeNotifiers()
	c.InitializeServices()
	c.InitializeHandlers()
	return c
}

func (c *Container) InitializeNotifiers() {
	c.Notifiers = []notifierInterface.Notifier{
		notifier.NewEmailNotifier(c.NotificationRepository),
	}
}

func (c *Container) InitializeRepositories() {
	// intialize shared repositories first
	c.ApiKeyRepository = api_key.NewApiKeyRepository(c.DB)

	c.AccountRepository = account.NewAccountRepository(c.DB, c.ApiKeyRepository)
	c.ConfigurationRepository = configuration.NewConfigurationRepository(c.DB)
	c.NotificationRepository = notification.NewNotificationRepository(c.DB)
}

func (c *Container) InitializeServices() {
	c.AccountService = account.NewAccountService(c.AccountRepository)
	c.ConfigurationService = configuration.NewConfigurationService(c.ConfigurationRepository)
	c.NotificationService = notification.NewNotificationService(c.Notifiers, c.NotificationRepository)
}

func (c *Container) InitializeHandlers() {
	c.AccountHandler = account.NewAccountHandler(c.AccountService)
	c.ConfigurationHandler = configuration.NewConfigurationHandler(c.ConfigurationService)
	c.NotificationHandler = notification.NewNotificationHandler(c.NotificationService)
}

