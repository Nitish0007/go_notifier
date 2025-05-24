package utils

import (
	"fmt"
	"log"
	"os"
	"context"

	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v2"
)

type dbConfig struct {
	Adapter	 string `yaml:"adapter"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type envConfig struct {
	Development dbConfig `yaml:"development"`
}

var DB *pgx.Conn
var dbConf dbConfig

// Public methods below
func ConnectDB() (*pgx.Conn, error) {
	setDbConfigs()
	dbURL := getDbURL()
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return	nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
		return nil, err
	}
	log.Println("Connected to database successfully")
	DB = conn
	return DB, nil
}

// Private methods below
func setDbConfigs() {	
	fileData, err := os.ReadFile("configs/database.yml")
	if(err != nil){
		log.Fatalf("Could not read config file: %v", err)
		panic(err)
	}

	var envConf envConfig
	err = yaml.Unmarshal(fileData, &envConf)
	if(err != nil){
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