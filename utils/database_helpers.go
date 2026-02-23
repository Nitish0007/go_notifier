package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Nitish0007/go_notifier/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/postgres"
	"gopkg.in/yaml.v2"
)

type dbConfig struct {
	Adapter  string `yaml:"adapter"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type envConfig struct {
	Development dbConfig `yaml:"development"`
}

var dbConf dbConfig


// Public methods below
func ConnectDB() (*gorm.DB, error) {
	setDbConfigs()
	dsn := getDSN()

	// to print sql queries in development environment synchronously
	customLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 1 * time.Second, // 1 second slow threshold for sql queries which are slow
			LogLevel: logger.Info, // log level for sql queries
			Colorful: true, // colorful output for sql queries
		},
	)

	// this db instance is used for Query building, model associations and ORM operations
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: customLogger})
	
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, err
	}
	
	// this is for connection pool management & health checks
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Unable to get underlying sql.DB: %v", err)
		return nil, err
	}
	
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
		return nil, err
	}

	// Configuring connection pool
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5* time.Minute)
	
	log.Println("Connected to database successfully with GORM")
	return db, nil
}

// Private methods below
func setDbConfigs() {
	fileData, err := os.ReadFile("configs/database.yml")
	if err != nil {
		log.Fatalf("Could not read config file: %v", err)
		panic(err)
	}

	var envConf envConfig
	err = yaml.Unmarshal(fileData, &envConf)
	if err != nil {
		log.Fatalf("Could not parse config file: %v", err)
		panic(err)
	}

	dbConf = envConf.Development
}

func getDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		dbConf.Host,
		dbConf.Username,
		dbConf.Password,
		dbConf.Database,
		dbConf.Port,
	)
}

// NOTE: Don't use this in production
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Account{},
		&models.Notification{},
		&models.NotificationBatch{},
		&models.NotificationBatchError{},
	)
}
