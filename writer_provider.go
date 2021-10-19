package polsvoice

import "io"

type WriterProvider interface {
	GetWriter() (io.Writer, string, func(), error)
}
