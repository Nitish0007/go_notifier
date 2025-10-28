// package utils

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"gopkg.in/yaml.v2"
// )

// type dbConfig struct {
// 	Adapter  string `yaml:"adapter"`
// 	Host     string `yaml:"host"`
// 	Port     int    `yaml:"port"`
// 	Username string `yaml:"username"`
// 	Password string `yaml:"password"`
// 	Database string `yaml:"database"`
// }

// type envConfig struct {
// 	Development dbConfig `yaml:"development"`
// }

// var dbConf dbConfig

// // var DB *pgxpool.Pool

// // Public methods below
// func ConnectDB() (*pgxpool.Pool, error) {
// 	setDbConfigs()
// 	dbURL := getDbURL()
// 	conn, err := pgxpool.New(context.Background(), dbURL)
// 	if err != nil {
// 		log.Fatalf("Unable to connect to database: %v", err)
// 		return nil, err
// 	}
// 	err = conn.Ping(context.Background())
// 	if err != nil {
// 		log.Fatalf("Unable to ping database: %v", err)
// 		return nil, err
// 	}
// 	log.Println("Connected to database successfully")
// 	return conn, nil
// }

// // Private methods below
// func setDbConfigs() {
// 	fileData, err := os.ReadFile("configs/database.yml")
// 	if err != nil {
// 		log.Fatalf("Could not read config file: %v", err)
// 		panic(err)
// 	}

// 	var envConf envConfig
// 	err = yaml.Unmarshal(fileData, &envConf)
// 	if err != nil {
// 		log.Fatalf("Could not parse config file: %v", err)
// 		panic(err)
// 	}

// 	dbConf = envConf.Development
// }

// func getDbURL() string {
// 	return fmt.Sprintf("%s://%s:%s@%s:%d/%s",
// 		dbConf.Adapter,
// 		dbConf.Username,
// 		dbConf.Password,
// 		dbConf.Host,
// 		dbConf.Port,
// 		dbConf.Database,
// 	)
// }


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
	
	// this db instance is used for Query building, model associations and ORM operations
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
	})
	
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
