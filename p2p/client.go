package p2p

import (
	"log"
	"net"

	"zc/blockchain"
	"zc/bullshit"
)

func StartClient(host string, port int, bc blockchain.Blockchain, p2pPort int) {
	addr := &net.TCPAddr{IP: net.ParseIP(host), Port: port}
	startClient(addr, bc, p2pPort, true, false)
}

func startClient(
	addr *net.TCPAddr,
	bc blockchain.Blockchain,
	p2pPort int,
	sendHandshake bool,
	reqBc bool,
) {
	log.Println("Client: dial to", addr)

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		bullshit.WarnIf(err)
		return
	}
	defer conn.Close()

	server := clientServer{
		conn: conn,
		r:    conn,
		w:    conn,
	}

	if sendHandshake {
		weight := bc.Weight()
		server.sendHandshake(p2pPort, weight)
	}

	if reqBc {
		server.reqChain()
	}

	for {
		resp, err := server.readResp()
		if err != nil {
			bullshit.WarnIf(err)
			return
		}

		switch resp {
		case respTypeBlock:
			block := server.readBlock()

			if !bc.AddBlock(block) {
				server.reqChain()
			} else {
				log.Println("Client: block added")
			}
		case respTypeChain:
			chain := server.readChain()

			if chain != nil && bc.Replace(chain) {
				log.Println("Client: chain replaced")
			}
		}
	}
}
