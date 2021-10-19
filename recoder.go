package polsvoice

import (
	"log"
	"os"
	"sync"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/youpy/go-wav"
)

const buffLen = 1500000
const numOfAudioChannels = 2
const samplingRate = 48000
const samplingBit = 16

type Recorder struct {
	samples        []wav.Sample
	writerProvider WriterProvider
	wg             sync.WaitGroup
}

func NewRecorder() *Recorder {
	return &Recorder{
		samples:        make([]wav.Sample, 0, buffLen),
		writerProvider: &FileWriterProvider{},
	}
}

func (r *Recorder) Record(vc *discordgo.VoiceConnection, finishChan chan os.Signal) error {
	rx := make(chan *discordgo.Packet, 2)

	go dgvoice.ReceivePCM(vc, rx)

	var p *discordgo.Packet
	var ok bool
	for {
		select {
		case <-finishChan:
			log.Println("finalizing...")

			err := r.writePCM()
			if err != nil {
				log.Println(err)
			}

			// MEMO should flush all remained packets in rx channel?

			r.wg.Wait()

			return nil
		case p, ok = <-rx:
			if !ok {
				return nil // TODO
			}
		}

		pcmLen := len(p.PCM)
		for i := 0; i < pcmLen; i += 2 {
			r.appendSample(wav.Sample{
				Values: [2]int{int(p.PCM[i]), int(p.PCM[i+1])},
			})

			if r.getSamplesLen() >= buffLen {
				err := r.writePCM()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func (r *Recorder) writePCM() error {
	defer r.clearSamples()

	w, identifier, closer, err := r.writerProvider.GetWriter()
	if err != nil {
		return err
	}
	log.Printf("write wav file; %s", identifier)

	samples := make([]wav.Sample, r.getSamplesLen())
	_ = copy(samples, r.samples)

	r.wg.Add(1)
	go func() {
		defer closer()
		defer r.wg.Done()

		err = wav.NewWriter(w, uint32(len(samples)), numOfAudioChannels, samplingRate, samplingBit).WriteSamples(samples)
		if err != nil {
			log.Println(err)
		}

		log.Printf("finished writing wav file; %s", identifier)
	}()

	return nil
}

func (r *Recorder) appendSample(sample wav.Sample) {
	r.samples = append(r.samples, sample)
}

func (r *Recorder) clearSamples() {
	r.samples = r.samples[:0]
}

func (r *Recorder) getSamplesLen() int {
	return len(r.samples)
}
