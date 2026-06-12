package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Options struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func New(o Options) (*sql.DB, error) {
	psqlConn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		o.Host, o.Port, o.User, o.Password, o.DBName,
	)

	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	log.Println("db.Stats", db.Stats())

	return db, nil
}
