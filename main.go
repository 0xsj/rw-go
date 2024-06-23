package main

import (
	"fmt"
	"os"

	"github.com/0xsj/rw-go/api"
	"github.com/0xsj/rw-go/config"
	db "github.com/0xsj/rw-go/db/sqlc"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("main running")
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "dev"
	}

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log := logger.Sugar()

	config := config.LoadConfig(env, ".")

	dbConn := db.Connect(config)
	defer db.Close(dbConn)

	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	store := db.New(dbConn) // Use the transaction object instead of dbConn
	server := api.NewServer(
		config,
		store,
		log,
	)

	server.MountHandlers()
	addr := fmt.Sprintf(":%s", config.Port)
	if err := server.Start(addr); err != nil {
		log.Fatal(err)
	}
}
