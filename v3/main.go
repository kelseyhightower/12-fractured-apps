package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	config Config
	db     *sql.DB
)

type Config struct {
	DataDir string `json:"datadir"`

	// Database settings.
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func main() {
	log.Println("Starting application...")

	// Load configuration settings.
	data, err := ioutil.ReadFile("/etc/config.json")

	switch {
	case os.IsNotExist(err):
		log.Println("Config file missing using defaults")
		config = Config{
			DataDir:  "/var/lib/data",
			Host:     "127.0.0.1",
			Port:     "3306",
			Database: "test",
		}
	case err == nil:
		if err := json.Unmarshal(data, &config); err != nil {
			log.Fatal(err)
		}
	default:
		log.Println(err)
	}

	log.Println("Overriding configuration from env vars.")
	if os.Getenv("APP_DATADIR") != "" {
		config.DataDir = os.Getenv("APP_DATADIR")
	}
	if os.Getenv("APP_HOST") != "" {
		config.Host = os.Getenv("APP_HOST")
	}
	if os.Getenv("APP_PORT") != "" {
		config.Port = os.Getenv("APP_PORT")
	}
	if os.Getenv("APP_USERNAME") != "" {
		config.Username = os.Getenv("APP_USERNAME")
	}
	if os.Getenv("APP_PASSWORD") != "" {
		config.Password = os.Getenv("APP_PASSWORD")
	}
	if os.Getenv("APP_DATABASE") != "" {
		config.Database = os.Getenv("APP_DATABASE")
	}

	// Use working directory.
	_, err = os.Stat(config.DataDir)
	if os.IsNotExist(err) {
		log.Println("Creating missing data directory", config.DataDir)
		err = os.MkdirAll(config.DataDir, 0755)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database.
	hostPort := net.JoinHostPort(config.Host, config.Port)
	log.Println("Connecting to database at", hostPort)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=30s",
		config.Username, config.Password, hostPort, config.Database)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
	}

	var dbError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		dbError = db.Ping()
		if dbError == nil {
			break
		}
		log.Println(dbError)
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	if dbError != nil {
		log.Fatal(dbError)
	}

	log.Println("Application started successfully.")
}
