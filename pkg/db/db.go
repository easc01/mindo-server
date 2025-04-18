package db

import (
	"database/sql"

	"github.com/ishantSikdar/mindo-server/internal/models"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
	_ "github.com/lib/pq"
)

var DB *sql.DB
var Queries *models.Queries

func InitDB(dsn string) error {
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if err := DB.Ping(); err != nil {
		return err
	}

	Queries = models.New(DB)
	logger.Log.Info("Database connected and Queries initialized")
	return nil
}
