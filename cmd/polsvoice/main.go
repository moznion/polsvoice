package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/moznion/polsvoice"
)

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := discord.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	vc, err := discord.ChannelVoiceJoin(os.Getenv("SERVER_ID"), os.Getenv("CHANNEL_ID"), true, false) // TODO
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := vc.Disconnect()
		if err != nil {
			log.Println(err)
		}
	}()

	finishChan := make(chan os.Signal, 1)
	signal.Notify(finishChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("start recording...")
	recorder := polsvoice.NewRecorder()
	err = recorder.Record(vc, finishChan)
	if err != nil {
		log.Println(err)
	}
}
