package db

import (
	"database/sql"
	"fmt"

	"github.com/ishantSikdar/mindo-server/internal/config"
	"github.com/ishantSikdar/mindo-server/internal/models"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
	_ "github.com/lib/pq"
)

var DB *sql.DB
var Queries *models.Queries

func InitDB() {
	var err error
	DB, err = sql.Open(
		"postgres",
		fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			config.GetConfig().DbUser,
			config.GetConfig().DbPassword,
			config.GetConfig().DbHost,
			config.GetConfig().DbPort,
			config.GetConfig().DbName,
		),
	)
	if err != nil {
		logger.Log.Error("DB connection failed", err)
		panic(err)
	}

	if err := DB.Ping(); err != nil {
		logger.Log.Errorf("db ping failure %s", err)
		panic(err)

	} else {
		logger.Log.Info("Database connected and Queries initialized")
	}

	Queries = models.New(DB)
}
