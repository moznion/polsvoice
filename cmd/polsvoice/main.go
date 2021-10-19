package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/moznion/polsvoice"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	err := polsvoice.Run(os.Getenv("BOT_TOKEN"), os.Getenv("SERVER_ID"), os.Getenv("CHANNEL_ID"))
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
