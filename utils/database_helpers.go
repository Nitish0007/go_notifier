package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
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
	Test        dbConfig `yaml:"test"`
	Production  dbConfig `yaml:"production"`
}

var dbConf dbConfig

// Public methods below
func ConnectDB(env string) (*gorm.DB, error) {
	// setDbConfigs(env)
	dsn := getDSN(env)

	// to print sql queries in development environment synchronously
	customLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 1 * time.Second, // 1 second slow threshold for sql queries which are slow
			LogLevel:      logger.Info,     // log level for sql queries
			Colorful:      true,            // colorful output for sql queries
		},
	)

	// this db instance is used for Query building, model associations and ORM operations
	db, err := gorm.Open(gormPostgres.Open(dsn), &gorm.Config{Logger: customLogger})

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
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	
	log.Println("Connected to database successfully with GORM")
	return db, nil
}

// SetupTestDB connects to the test database, runs migrations, and returns a GORM DB instance
// This method is intended for use in integration tests only
func SetupTestDB() (*gorm.DB, error) {
	// Get test database DSN
	dsn := getDSN("test")

	// Convert DSN to database URL format for migrate
	dbConf := getDBConfigs("test")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConf.Username,
		dbConf.Password,
		dbConf.Host,
		dbConf.Port,
		dbConf.Database,
	)

	// Open database connection using database/sql (required by migrate)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create postgres driver instance for migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Get absolute path to migrations directory
	migrationsDir, err := filepath.Abs("db/migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path to migrations: %w", err)
	}

	// Verify migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("migrations directory does not exist: %s", migrationsDir)
	}

	// Create file source (reads migration files from directory)
	source, err := (&file.File{}).Open("file://" + migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open migrations source: %w", err)
	}
	defer source.Close()

	// Create migrate instance
	m, err := migrate.NewWithInstance(
		"file",
		source,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run all migrations up
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		// ErrNoChange means migrations are already applied, which is fine
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Test database migrations completed successfully")

	// Now connect with GORM for use in tests
	gormDB, err := gorm.Open(gormPostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect with GORM: %w", err)
	}

	// Verify GORM connection
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database with GORM: %w", err)
	}

	log.Println("Test database connected successfully with GORM")
	return gormDB, nil
}


// Private methods below

func getDBConfigs(env string) dbConfig {
	if env == "" {
		env = "development"
	}

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
	
	switch env {
	case "production":
		return envConf.Production
	case "test":
		return envConf.Test
	default:
		return envConf.Development
	}
}

func getDSN(env string) string {
	dbConf := getDBConfigs(env)
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		dbConf.Host,
		dbConf.Username,
		dbConf.Password,
		dbConf.Database,
		dbConf.Port,
	)
}

