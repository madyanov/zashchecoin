package p2p

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"

	"zc/blockchain"
)

const (
	idleTimeout = 10 * time.Minute
)

type serverClient struct {
	conn net.Conn
	id   int
	r    io.Reader
	w    io.Writer
}

func newServerClient(
	conn net.Conn,
	id int,
	r io.Reader,
	w io.Writer,
) serverClient {
	client := serverClient{
		conn: conn,
		id:   id,
		r:    conn,
		w:    conn,
	}

	client.updateDeadline()
	return client
}

func (c serverClient) sendBlock(block blockchain.Block) {
	log.Println("Server: send block")
	c.updateDeadline()
	resp := respBlock{block: block}
	resp.encode(c.w)
}

func (c serverClient) sendBlockchain(bc blockchain.Blockchain) {
	log.Println("Server: send blockchain")
	c.updateDeadline()
	resp := respChain{blocks: bc.AllBlocks()}
	resp.encode(c.w)
}

func (c serverClient) readReq() (byte, error) {
	log.Println("Server: read request type")
	c.updateDeadline()
	var req byte
	err := binary.Read(c.r, binary.BigEndian, &req)
	return req, err
}

func (c serverClient) readHandshake() reqHandshake {
	log.Println("Server: read handshake")
	c.updateDeadline()
	req := reqHandshake{}
	req.decode(c.r)
	return req
}

func (c serverClient) updateDeadline() {
	deadline := time.Now().Add(idleTimeout)
	c.conn.SetDeadline(deadline)
}
