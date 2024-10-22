package storage

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func NewMysqlClient() {
	c := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		DBName:               os.Getenv("MYSQL_DATABASE"),
		Addr:                 "localhost:3306",
		Net:                  "tcp",
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = db.Close()
		log.Println("db connection closed")
	}()

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}
}
