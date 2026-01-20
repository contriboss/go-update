package binarydist

import (
	"io"

	"github.com/dsnet/compress/bzip2"
)

type bzip2Writer struct {
	w *bzip2.Writer
}

func (w bzip2Writer) Write(b []byte) (int, error) {
	return w.w.Write(b)
}

func (w bzip2Writer) Close() error {
	return w.w.Close()
}

// newBzip2Writer creates a new bzip2 writer using pure Go implementation.
func newBzip2Writer(w io.Writer) (io.WriteCloser, error) {
	bw, err := bzip2.NewWriter(w, &bzip2.WriterConfig{Level: bzip2.DefaultCompression})
	if err != nil {
		return nil, err
	}
	return bzip2Writer{w: bw}, nil
}
