package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
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

// var DB *pgxpool.Pool

// Public methods below
func ConnectDB() (*pgxpool.Pool, error) {
	setDbConfigs()
	dbURL := getDbURL()
	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
		return nil, err
	}
	log.Println("Connected to database successfully")
	return conn, nil
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

func getDbURL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		dbConf.Adapter,
		dbConf.Username,
		dbConf.Password,
		dbConf.Host,
		dbConf.Port,
		dbConf.Database,
	)
}
