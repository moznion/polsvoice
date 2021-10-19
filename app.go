package polsvoice

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Run is an entry point of the app.
func Run(botToken string, serverID string, channelID string, filePrefix string) error {
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
			log.Error().Err(err).Msg("failed to close discord instance")
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

	log.Info().Msg("standby recording in speaker mute")
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
			log.Error().Err(err).Msg("failed to disconnect from voice channel in mute (i.e. initial state)")
		}
	}()

	select {
	case <-finishChan:
		log.Info().Msg("finished app")
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
			log.Error().Err(err).Msg("failed to disconnect from voice channel")
		}
	}()

	log.Info().Msg("start recording...")
	recorder := NewRecorder(filePrefix)
	err = recorder.Record(vc, finishChan)
	if err != nil {
		return err
	}

	return nil
}
