package polsvoice

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Run is an entry point of the app.
func Run(botToken string, serverID string, channelID string) error {
	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return err
	}

	err = discord.Open()
	if err != nil {
		return err
	}
	defer func() {
		err := discord.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	vc, err := discord.ChannelVoiceJoin(serverID, channelID, true, false) // TODO
	if err != nil {
		return err
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
	recorder := NewRecorder()
	err = recorder.Record(vc, finishChan)
	if err != nil {
		return err
	}

	return nil
}
