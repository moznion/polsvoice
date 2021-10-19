package main

import (
	"log"
	"os"

	"github.com/moznion/polsvoice"
)

func main() {
	err := polsvoice.Run(os.Getenv("BOT_TOKEN"), os.Getenv("SERVER_ID"), os.Getenv("CHANNEL_ID"))
	if err != nil {
		log.Fatal(err)
	}
}
