package main

import (
	"coding-challenge-go/pkg/api"
	database "coding-challenge-go/pkg/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	// TODO: make db on error exit more explicit
	dbNew := database.InitConnection("mysql", "user:password@tcp(db:3306)/product")

	defer dbNew.Close()

	engine, err := api.CreateAPIEngine(dbNew)

	if err != nil {
		log.Error().Err(err).Msg("Fail to create server")
		return
	}

	log.Info().Msg("Start server")
	log.Fatal().Err(engine.Run(os.Getenv("LISTEN"))).Msg("Fail to listen and serve")
}
