package p2p

import (
	"log"
	"net"
	"strconv"
	"sync"

	"zc/blockchain"
	"zc/bullshit"
)

type server struct {
	clients    map[int]serverClient
	clientId   int
	clientsMtx sync.RWMutex
}

func StartServer(port int, bc *blockchain.Blockchain, blockChan <-chan blockchain.Block) {
	log.Printf("Server: listen on %d", port)

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	bullshit.FailIf(err)
	defer ln.Close()

	server := &server{}
	server.clients = make(map[int]serverClient)

	for {
		conn, err := ln.Accept()
		bullshit.WarnIf(err)

		go server.handleConn(conn, bc, port)
		go server.broadcastBlocks(blockChan)
	}
}

func (s *server) handleConn(conn net.Conn, bc *blockchain.Blockchain, port int) {
	defer conn.Close()

	client := s.addClient(conn)
	defer s.removeClient(client)

	for {
		req, err := client.readReq()
		if err != nil {
			bullshit.WarnIf(err)
			return
		}

		switch req {
		case reqTypeHandshake:
			handshake := client.readHandshake()

			if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
				addr.Port = int(handshake.port)
				go startClient(addr, bc, port, false)
			}
		case reqTypeChain:
			client.sendBlockchain(bc)
		}
	}
}

func (s *server) broadcastBlocks(blockChan <-chan blockchain.Block) {
	for block := range blockChan {
		s.broadcastBlock(block)
	}
}

func (s *server) broadcastBlock(block blockchain.Block) {
	s.clientsMtx.RLock()
	defer s.clientsMtx.RUnlock()

	for _, client := range s.clients {
		client.sendBlock(block)
	}
}

func (s *server) addClient(conn net.Conn) serverClient {
	log.Println("Server: add client")

	s.clientsMtx.Lock()
	defer s.clientsMtx.Unlock()
	defer func() { s.clientId++ }()

	client := newServerClient(conn, s.clientId, conn, conn)
	s.clients[s.clientId] = client
	return client
}

func (s *server) removeClient(client serverClient) {
	log.Println("Server: remove client")

	s.clientsMtx.Lock()
	defer s.clientsMtx.Unlock()

	delete(s.clients, client.id)
}
