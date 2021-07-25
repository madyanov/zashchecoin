package p2p

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"zc/blockchain"
	"zc/bullshit"
)

type clientServer struct {
	conn net.Conn
	r    io.Reader
	w    io.Writer
}

func (c clientServer) sendHandshake(port int, weight int) {
	log.Println("Client: send handshake")

	req := reqHandshake{
		port:   int16(port),
		weight: int64(weight),
	}

	req.encode(c.w)
}

func (c clientServer) reqChain() {
	log.Println("Client: request chain")
	err := binary.Write(c.w, binary.BigEndian, reqTypeChain)
	bullshit.WarnIf(err)
}

func (c clientServer) readResp() (byte, error) {
	log.Println("Client: read response type")
	var resp byte
	err := binary.Read(c.r, binary.BigEndian, &resp)
	return resp, err
}

func (c clientServer) readBlock() blockchain.Block {
	log.Println("Client: read block")
	resp := respBlock{}
	resp.decode(c.r)
	return resp.block
}

func (c clientServer) readChain() *blockchain.MemBlockchain {
	log.Println("Client: read chain")

	resp := respChain{}
	resp.decode(c.r)

	bc, err := blockchain.NewMemBlockchain(resp.blocks)
	bullshit.WarnIf(err)

	return bc
}
