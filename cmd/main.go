package main

import (
	"github.com/easc01/mindo-server/internal/handlers"
	"github.com/easc01/mindo-server/pkg/db"
)

func main() {
	db.InitDB()
	handlers.InitREST()
}
