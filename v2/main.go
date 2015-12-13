package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

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
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	// Use working directory.
	_, err = os.Stat(config.DataDir)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database.
	hostPort := net.JoinHostPort(config.Host, config.Port)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=30s",
		config.Username, config.Password, hostPort, config.Database)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}
