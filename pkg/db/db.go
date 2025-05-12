package db

import (
	"database/sql"

	"github.com/easc01/mindo-server/internal/config"
	"github.com/easc01/mindo-server/internal/models"
	"github.com/easc01/mindo-server/pkg/logger"
	_ "github.com/lib/pq"
)

var DB *sql.DB
var Queries *models.Queries

func InitDB() {
	var err error
	DB, err = sql.Open("postgres", config.GetConfig().DbConnectionUri)
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
