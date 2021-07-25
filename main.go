package main

import (
	"bufio"
	"flag"
	"os"
	"strconv"
	"strings"

	"zc/api"
	"zc/blockchain"
	"zc/bullshit"
	"zc/db"
	"zc/p2p"
)

func main() {
	apiPort := flag.Int("api-port", 8080, "HTTP API port")
	p2pPort := flag.Int("p2p-port", 8081, "P2P server port")
	peersPath := flag.String("peers", "peers.txt", "Path to the peers file")
	bcPath := flag.String("bc", "bc.dat", "Path to the blockchain file")
	flag.Parse()

	bc := db.NewDiskBlockchain(*bcPath)

	blockChan := make(chan blockchain.Block)

	go api.StartServer(*apiPort, bc, blockChan)
	go p2p.StartServer(*p2pPort, bc, blockChan)

	connectToPeers(*peersPath, bc, *p2pPort)

	select {}
}

func connectToPeers(peersPath string, bc blockchain.Blockchain, p2pPort int) {
	f, err := os.Open(peersPath)
	if err != nil {
		bullshit.WarnIf(err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		peer := scanner.Text()
		if len(peer) == 0 {
			continue
		}

		parts := strings.Split(peer, ":")

		host := parts[0]
		port, err := strconv.Atoi(parts[1])
		bullshit.WarnIf(err)

		go p2p.StartClient(host, port, bc, p2pPort)
	}
}
