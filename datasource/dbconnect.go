package datasource

import (
	"database/sql"
	"streamapi/utils/logger"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Client *sql.DB
)

func init() {

	logger.Info.Println("Setting up database")
	var err error
	Client, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/streamdb")

	if err != nil {
		logger.Error.Println("Databse Setup error")
	}

	logger.Info.Println("Database Successfully Connected")

}
