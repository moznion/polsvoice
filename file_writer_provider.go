package polsvoice

import (
	"fmt"
	"io"
	"log"
	"os"
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
			log.Println(err)
		}
	}, nil

}
