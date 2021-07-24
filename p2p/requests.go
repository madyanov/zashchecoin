package p2p

import (
	"encoding/binary"
	"encoding/gob"
	"io"

	"zc/bullshit"
)

const (
	reqTypeHandshake = byte(iota)
	reqTypeChain
)

// Handshake
type reqHandshake struct {
	port int16
}

func (req *reqHandshake) encode(w io.Writer) {
	err := binary.Write(w, binary.BigEndian, reqTypeHandshake)
	bullshit.WarnIf(err)

	enc := gob.NewEncoder(w)
	err = enc.Encode(req.port)
	bullshit.WarnIf(err)
}

func (req *reqHandshake) decode(r io.Reader) {
	dec := gob.NewDecoder(r)
	err := dec.Decode(&req.port)
	bullshit.WarnIf(err)
}
