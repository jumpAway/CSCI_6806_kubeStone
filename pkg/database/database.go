package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"kubeStone/m/v2/pkg/config"
	"log"
	"strconv"
)

// InitDB is a function that initializes a connection to a MySQL database using the provided configuration.
func InitDB(config config.Config) (*sql.DB, error) {
	dsn := config.Database.Username + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + strconv.Itoa(config.Database.Port) + ")/" + config.Database.DatabaseName
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Connect to mysql error")
		return nil, err
	}
	return db, nil
}
