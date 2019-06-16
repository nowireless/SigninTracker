package database

import (
	"fmt"
	"signin3/models"

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

func (db *Database) Create(obj models.Model) error {
	switch v := obj.(type) {
	case *models.Person:
		return db.CreatePerson(v)
	}

	return fmt.Errorf("Unsupported type: %T", obj)
}

func (db *Database) Get(objType interface{}, databaseID int) (interface{}, error) {
	switch v := objType.(type) {
	case models.Person:
		return db.GetPerson(databaseID)
	default:
		return nil, fmt.Errorf("Unsupported type: %T", v)

	}
}

func (db *Database) Update(obj interface{}) error {
	switch v := obj.(type) {
	case models.Person:
		return db.UpdatePerson(v)
	default:
		return fmt.Errorf("Unsupported type: %T", v)
	}
}

func (db *Database) Delete(obj interface{}) error {
	switch v := obj.(type) {
	case *models.Person:
		return db.DeletePerson(v)
	default:
		return fmt.Errorf("Unsupported type: %T", v)
	}
}
