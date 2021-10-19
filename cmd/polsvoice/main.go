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
	var outputStr string
	var filePrefix string

	flag.StringVar(&botToken, "bot-token", "", "[mandatory] secret token of the Discord bot")
	flag.StringVar(&serverID, "server-id", "", "[mandatory] server ID to join in")
	flag.StringVar(&channelID, "channel-id", "", "[mandatory] voice channel ID to join in")
	flag.StringVar(&outputStr, "out", "file", "output destination; this parameter must be \"file\" or \"stdout\"")
	flag.StringVar(&filePrefix, "file-prefix", "", "output file prefix. if this value is \"test\", this will make \"test-${sequenceNo}.wav\". when you specify \"file\" as the output destination, this parameter is mandatory")

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

	var output polsvoice.OutType
	if outputStr == polsvoice.File {
		output = polsvoice.File
	} else if outputStr == polsvoice.Stdout {
		log.Fatal().Msg("TODO: not implemented yet")
	} else {
		log.Fatal().Msg("invalid --out parameter has come")
	}

	if output == polsvoice.File && filePrefix == "" {
		log.Fatal().Msg("when --out parameter is \"file\", \"--file-prefix\" must be specified")
	}

	err := polsvoice.Run(botToken, serverID, channelID, filePrefix)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
