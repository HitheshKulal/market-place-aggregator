package main

import (
	"go-backend-api/internal/bootstrap"
	"go-backend-api/internal/db"
	"go-backend-api/internal/logger"
)

func main() {
	bootstrap.LoadEnv()
	cfg := bootstrap.InitConfig()
	database := bootstrap.InitDB(cfg)
	if err := db.RunMigrations(database); err != nil {
		logger.Fatal("failed to run migrations: %w")
	}
	repos := bootstrap.InitRepositories(database)

	services := bootstrap.InitServices(
		repos,
	)
	router := bootstrap.InitRouter(services, cfg)
	server := bootstrap.InitServer(router, cfg)

	bootstrap.StartServer(server, cfg)
}
