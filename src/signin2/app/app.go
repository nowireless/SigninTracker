package app

import (
	"signin2/database"
)

type Config struct {
	Database database.Config
}

type App struct {
	DB     *database.Database
	Config Config
}

func NewApp(config Config) (*App, error) {
	var err error
	app := App{Config: config}
	app.DB, err = database.Connect(config.Database)
	if err != nil {
		return nil, err
	}

	return &app, nil
}
