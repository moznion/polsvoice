package polsvoice

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type MentionHandler struct {
	StartRecChan      chan interface{}
	DisconnectionChan chan interface{}
}

var recMsgRe = regexp.MustCompile("\\s+rec$")
var finishMsgRe = regexp.MustCompile("\\s+finish$")

func (h *MentionHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	selfUserID := s.State.User.ID

	if m.Author.ID == selfUserID { // ignore the message comes from itself
		return
	}

	isMentioned := false
	for _, mention := range m.Mentions {
		if mention.ID == selfUserID {
			isMentioned = true
			break
		}
	}

	if !isMentioned {
		// do nothing when the message hasn't mentioned the bot
		return
	}

	if recMsgRe.MatchString(m.Content) {
		h.handleRecMessage(s, m)
		return
	}

	if finishMsgRe.MatchString(m.Content) {
		h.handleFinishMessage(s, m)
		return
	}

	h.handleUnknownMessage(s, m)
}

func (h *MentionHandler) handleRecMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	select {
	case h.StartRecChan <- struct{}{}:
	default:
		return
	}

	_, err := s.ChannelMessageSend(m.ChannelID, "Okay, let's start recording")
	if err != nil {
		log.Error().Err(err).Msg("failed to response to rec msg")
	}
}

func (h *MentionHandler) handleFinishMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Bye...")
	if err != nil {
		log.Error().Err(err).Msg("failed to response to finish msg")
	}

	select {
	case h.DisconnectionChan <- struct{}{}:
	default:
	}
}

func (h *MentionHandler) handleUnknownMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "I don't know what you mean.")
	if err != nil {
		log.Error().Err(err).Msg("failed to response to unknown msg")
	}
}
