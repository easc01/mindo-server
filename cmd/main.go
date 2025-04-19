package main

import (
	"github.com/ishantSikdar/mindo-server/internal/handlers"
	"github.com/ishantSikdar/mindo-server/pkg/db"
)

func main() {
	db.InitDB()
	handlers.InitREST()
}
