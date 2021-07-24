package api

import (
	"log"
	"net/http"
	"strconv"

	"zc/blockchain"
	"zc/bullshit"
)

func StartServer(port int, bc *blockchain.Blockchain, blockChan chan<- blockchain.Block) {
	http.HandleFunc("/mine", mineBlock(bc, blockChan))

	log.Printf("API Server: listen on %d\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	bullshit.FailIf(err)
}

func mineBlock(bc *blockchain.Blockchain, blockChan chan<- blockchain.Block) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Println("API Server: mine block")

		data := req.URL.Query().Get("data")
		block := bc.MineBlock([]byte(data))

		if bc.AddBlock(block) {
			blockChan <- block
		}
	}
}
