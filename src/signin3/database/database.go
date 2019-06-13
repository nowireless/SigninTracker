package database

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	log "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

type Database struct {
	DB *sqlx.DB
}

func Connect(config Config) (*Database, error) {
	// See: https://github.com/jackc/pgx/blob/master/stdlib/sql.go
	connectionStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d database=%s sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	log.Info("Connection String: ", connectionStr)

	db, err := sqlx.Connect("pgx", connectionStr)
	if err != nil {
		return nil, err
	}
	log.Info("Pinging Database")
	db.Ping()

	return &Database{DB: db}, nil
}
