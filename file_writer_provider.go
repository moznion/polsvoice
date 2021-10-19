package polsvoice

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

type FileWriterProvider struct {
	fileSeqNum uint
}

func (p *FileWriterProvider) GetWriter() (io.Writer, string, func(), error) {
	identifier := fmt.Sprintf("test-%09d.wav", p.fileSeqNum)
	f, err := os.Create(identifier)
	if err != nil {
		return nil, "", func() {}, err
	}
	p.fileSeqNum++

	return f, identifier, func() {
		err := f.Close()
		if err != nil {
			log.Error().Err(err).Str("file_identifier", identifier).Msg("failed to close the file")
		}
	}, nil

}
