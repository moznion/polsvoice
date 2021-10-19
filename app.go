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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	disconnectionChan := make(chan interface{}, 1)

	finishChan := make(chan interface{}, 1)
	go func() {
		select {
		case <-disconnectionChan:
			finishChan <- struct{}{}
		case <-sigChan:
			finishChan <- struct{}{}
		}
	}()

	startRecChan := make(chan interface{}, 1)
	mentionHandler := &MentionHandler{
		StartRecChan:      startRecChan,
		DisconnectionChan: disconnectionChan,
	}
	discord.AddHandler(mentionHandler.Handle)

	log.Println("standby recording in speaker mute")
	alreadyDisconnectedChan := make(chan interface{}, 1)
	initialVC, err := discord.ChannelVoiceJoin(serverID, channelID, true, true)
	if err != nil {
		return err
	}
	defer func() {
		select {
		case <-alreadyDisconnectedChan:
			return
		default:
		}

		err := initialVC.Disconnect()
		if err != nil {
			log.Println(err)
		}
	}()

	select {
	case <-finishChan:
		log.Println("finished")
		return nil
	case <-startRecChan:
		// fall through
	}

	// reconnect to listen the voice channel
	err = initialVC.Disconnect()
	if err != nil {
		return err
	}
	vc, err := discord.ChannelVoiceJoin(serverID, channelID, true, false)
	if err != nil {
		return err
	}
	defer func() {
		err = initialVC.Disconnect()
		if err != nil {
			log.Println(err)
		}
	}()

	log.Println("start recording...")
	recorder := NewRecorder()
	err = recorder.Record(vc, finishChan)
	if err != nil {
		return err
	}

	return nil
}
