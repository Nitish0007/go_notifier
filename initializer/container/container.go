package container

import (
	"gorm.io/gorm"

	"github.com/Nitish0007/go_notifier/internal/features/account"
	"github.com/Nitish0007/go_notifier/internal/features/apiKey"
	"github.com/Nitish0007/go_notifier/internal/features/configuration"
	"github.com/Nitish0007/go_notifier/internal/features/contact"
	"github.com/Nitish0007/go_notifier/internal/features/content"
	"github.com/Nitish0007/go_notifier/internal/features/emailcontact"
	"github.com/Nitish0007/go_notifier/internal/features/emailnotification"
	"github.com/Nitish0007/go_notifier/internal/features/emailnotificationlist"
	"github.com/Nitish0007/go_notifier/internal/features/list"
	"github.com/Nitish0007/go_notifier/internal/features/listsubscription"
	// "github.com/Nitish0007/go_notifier/internal/lib/notifier"
	// notifierInterface "github.com/Nitish0007/go_notifier/internal/shared/interfaces/notifier"
)

type Container struct {
	DB *gorm.DB
	// NOTE: keep in mind that dependencies should be intialized first
	// repositories are the lowest level of dependencies and should be intialized first

	// Repositories
	ApiKeyRepository                *apiKey.ApiKeyRepository
	AccountRepository               *account.AccountRepository
	ConfigurationRepository         *configuration.ConfigurationRepository
	EmailNotificationRepository     *emailnotification.EmailNotificationRepository
	EmailNotificationListRepository *emailnotificationlist.EmailNotificationListRepository
	ContactRepository               *contact.ContactRepository
	EmailContactRepository          *emailcontact.EmailContactRepository
	ConentRepository                *content.ContentRepository
	ListSubscriptionRepository      *listsubscription.ListSubscriptionRepository
	ListRepository                  *list.ListRepository

	// Notifiers
	// Notifiers []notifierInterface.Notifier

	// Services
	AccountService           *account.AccountService
	ConfigurationService     *configuration.ConfigurationService
	EmailNotificationService *emailnotification.EmailNotificationService
	ContactService           *contact.ContactService
	ListService              *list.ListService
	ContentService           *content.ContentService

	// Handlers
	AccountHandler           *account.AccountHandler
	ConfigurationHandler     *configuration.ConfigurationHandler
	EmailNotificationHandler *emailnotification.EmailNotificationHandler
	ContactHandler           *contact.ContactHandler
	ContentHandler           *content.ContentHandler
	ListHandler              *list.ListHandler
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

	c.AccountRepository               = account.NewAccountRepository(c.DB, c.ApiKeyRepository)
	c.ConfigurationRepository         = configuration.NewConfigurationRepository(c.DB)
	c.ContactRepository               = contact.NewContactRepository(c.DB, c.EmailContactRepository)
	c.EmailContactRepository          = emailcontact.NewEmailContactRepository(c.DB)
	c.EmailNotificationListRepository = emailnotificationlist.NewEmailNotificationListRepository(c.DB)
	c.EmailNotificationRepository     = emailnotification.NewEmailNotificationRepository(c.DB, c.EmailNotificationListRepository)
	c.ListSubscriptionRepository      = listsubscription.NewListSubscriptionRepository(c.DB)
	c.ConentRepository                = content.NewContentRepository(c.DB)
	c.ListRepository                  = list.NewListRepository(c.DB, c.ListSubscriptionRepository, c.ContactRepository)
}

func (c *Container) InitializeServices() {
	c.AccountService           = account.NewAccountService(c.AccountRepository)
	c.ConfigurationService     = configuration.NewConfigurationService(c.ConfigurationRepository)
	c.ContactService           = contact.NewContactService(c.ContactRepository)
	c.EmailNotificationService = emailnotification.NewEmailNotificationService(c.EmailNotificationRepository)
	c.ContentService           = content.NewContentService(c.ConentRepository)
	c.ListService              = list.NewListService(c.ListRepository)
}

func (c *Container) InitializeHandlers() {
	c.AccountHandler           = account.NewAccountHandler(c.AccountService)
	c.ConfigurationHandler     = configuration.NewConfigurationHandler(c.ConfigurationService)
	c.ContactHandler           = contact.NewContactHandler(c.ContactService)
	c.EmailNotificationHandler = emailnotification.NewEmailNotificationHandler(c.EmailNotificationService)
	c.ListHandler              = list.NewListHandler(c.ListService)
	c.ContentHandler           = content.NewContentHandler(c.ContentService)
}
