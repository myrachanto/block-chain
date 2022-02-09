package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/myrachanto/blockchain/block"
	"github.com/myrachanto/blockchain/wallet"
)

var cache map[string]*block.Blockchain = make(map[string]*block.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}
func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}
func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("Error: Invalid HTTp method")
	}
}
func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = block.NewBlockchain(minersWallet.BlockchainAddress(), bcs.Port())
		cache["blockchain"] = bc
		log.Printf("private key %v", minersWallet.PrivateKeyStr())
		log.Printf("public key %v", minersWallet.PublicKeyStr())
		log.Printf("Blockchain address %v", minersWallet.BlockchainAddress())
	}
	return bc
}
func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", bcs.GetChain)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.Port())), nil))
}
