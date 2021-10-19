package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/moznion/polsvoice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	var botToken string
	var serverID string
	var channelID string

	flag.StringVar(&botToken, "bot-token", "", "[mandatory] secret token of the Discord bot")
	flag.StringVar(&serverID, "server-id", "", "[mandatory] server ID to join in")
	flag.StringVar(&channelID, "channel-id", "", "[mandatory] voice channel ID to join in")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, `polsvoice: A Discord bot to record the sound in voice chat.

Usage of %s:
   %s [OPTIONS]
Options
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if botToken == "" {
		flag.Usage()
		log.Fatal().Msg("--bot-token is a mandatory parameter")
	}
	if serverID == "" {
		flag.Usage()
		log.Fatal().Msg("--server-id is a mandatory parameter")
	}
	if channelID == "" {
		flag.Usage()
		log.Fatal().Msg("--channel-id is a mandatory parameter")
	}

	err := polsvoice.Run(botToken, serverID, channelID)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
