package storage

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Storage struct {
	db     *sql.DB
	logger *slog.Logger
	// TickChan       chan *models.Tick
	// BidAskTickChan chan *models.BidAskTick
}

func NewMysqlClient(logger *slog.Logger) (*Storage, error) {
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

	err = db.Ping()

	if err != nil {
		return &Storage{
			db: nil,
		}, err
	}

	return &Storage{
		db:     db,
		logger: logger,
		// TickChan:       make(chan *models.Tick, 1), // buffer channel with size 1
		// BidAskTickChan: make(chan *models.BidAskTick, 1),
	}, nil
}

func (s *Storage) CreateTableKLine() {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS ticks (
		id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, 
		coin VARCHAR(10) NOT NULL, 
		price DECIMAL(18,8) NOT NULL, 
		dt TIMESTAMP
	)`)

	if err != nil {
		log.Fatal("error creating kLine table", err)
	}
}

func (s *Storage) Close() {
	s.db.Close()
	log.Println("db connection closed")
}
