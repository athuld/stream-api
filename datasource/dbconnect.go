package datasource

import (
	"database/sql"
	"os"
	"streamapi/utils/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	Client *sql.DB
)

func init() {

	godotenv.Load()
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	dbIP := os.Getenv("DB_IP")
	logger.Info.Println("Setting up database")
	var err error
	Client, err = sql.Open("mysql", username+":"+password+"@tcp("+dbIP+")/streamdb")

	if err != nil {
		logger.Error.Println("Databse Setup error")
	}

	logger.Info.Println("Database Successfully Connected")

}
