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
	"zc/p2p"
)

func main() {
	apiPort := flag.Int("api-port", 8080, "HTTP API port")
	srvPort := flag.Int("srv-port", 8081, "P2P server port")
	peersPath := flag.String("peers", "peers.txt", "Path to peers file")
	flag.Parse()

	bc := blockchain.NewBlockchain([]blockchain.Block{})

	blockChan := make(chan blockchain.Block)

	go api.StartServer(*apiPort, bc, blockChan)
	go p2p.StartServer(*srvPort, bc, blockChan)

	connectToPeers(*peersPath, bc, *srvPort)

	select {}
}

func connectToPeers(peersPath string, bc *blockchain.Blockchain, srvPort int) {
	file, err := os.Open(peersPath)
	if err != nil {
		bullshit.WarnIf(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		peer := scanner.Text()
		if len(peer) == 0 {
			continue
		}

		parts := strings.Split(peer, ":")

		host := parts[0]
		port, err := strconv.Atoi(parts[1])
		bullshit.WarnIf(err)

		go p2p.StartClient(host, port, bc, srvPort)
	}
}
