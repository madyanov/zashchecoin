package encoding

import (
	"bytes"
	"encoding/gob"
	"io"
)

// https://stackoverflow.com/a/60827466

func NewEncoderWithoutHeader(w io.Writer, dummy interface{}) (*gob.Encoder, *bytes.Buffer, error) {
	rw := &readerWriter{}
	buf := &bytes.Buffer{}
	rw.w = buf

	enc := gob.NewEncoder(rw)
	if err := enc.Encode(dummy); err != nil {
		return nil, nil, err
	}

	rw.w = w
	return enc, buf, nil
}

func NewDecoderWithoutHeader(r io.Reader, buf *bytes.Buffer, dummy interface{}) (*gob.Decoder, error) {
	rw := &readerWriter{}
	rw.r = buf

	dec := gob.NewDecoder(rw)
	if err := dec.Decode(dummy); err != nil {
		return nil, err
	}

	rw.r = r
	return dec, nil
}

type readerWriter struct {
	r io.Reader
	w io.Writer
}

func (rw *readerWriter) Read(p []byte) (int, error) {
	return rw.r.Read(p)
}

func (rw *readerWriter) Write(p []byte) (int, error) {
	return rw.w.Write(p)
}
