package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/youpy/go-wav"
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

	seqNo := 0
	fileSeqNo := 0

	buffLen := 1500000

	recv := make(chan *discordgo.Packet, 1024)
	samples := make([]wav.Sample, 0, buffLen)
	go dgvoice.ReceivePCM(vc, recv)
	for {
		p, ok := <-recv
		if !ok {
			os.Exit(1) // TODO
		}

		l := len(p.PCM)

		for i := 0; i < l; i += 2 {
			samples = append(samples, wav.Sample{
				Values: [2]int{int(p.PCM[i]), int(p.PCM[i+1])},
			})

			seqNo++
			if seqNo >= buffLen {
				func() {
					f, err := os.Create(fmt.Sprintf("test-%d.wav", fileSeqNo))
					if err != nil {
						log.Println(err)
						return
					}
					defer func() {
						f.Close()
					}()
					fileSeqNo++

					w := wav.NewWriter(f, uint32(len(samples)), 2, 48000, 16)
					err = w.WriteSamples(samples)
					if err != nil {
						log.Println(err)
					}
				}()

				seqNo = 0
				samples = make([]wav.Sample, 0, buffLen)
			}
		}
	}
}
