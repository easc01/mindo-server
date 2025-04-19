package db

import (
	"database/sql"

	"github.com/ishantSikdar/mindo-server/internal/models"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
	_ "github.com/lib/pq"
)

var DB *sql.DB
var Queries *models.Queries

func InitDB() {
	var err error
	DB, err = sql.Open("postgres", "postgresql://postgres:6515@localhost:5432/mindo?sslmode=disable")
	if err != nil {
		logger.Log.Error("DB connection failed", err)

	}

	if err := DB.Ping(); err != nil {
		logger.Log.Error("DB ping failure", err)

	}

	Queries = models.New(DB)
	logger.Log.Info("Database connected and Queries initialized")
}
