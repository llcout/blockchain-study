package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/lucasacoutinho/astroboi/block"
	"github.com/lucasacoutinho/astroboi/wallet"
)

var cache map[string]*block.BlockChain = make(map[string]*block.BlockChain)

type Server struct {
	port uint16
}

func NewServer(port uint16) *Server {
	return &Server{port: port}
}

func (s *Server) Port() uint16 {
	return s.port
}

func (s *Server) GetBlockchain() *block.BlockChain {
	bc, ok := cache["blockchain"]
	if !ok {
		wa := wallet.NewWallet()
		bc := block.NewBlockChain(wa.Address(), s.Port())
		cache["blockchain"] = bc
	}

	return bc
}
func (s *Server) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-type", "application/json")
		bc := s.GetBlockchain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: invalid method")
	}
}

func (s *Server) Run() {
	http.HandleFunc("/", s.GetChain)

	address := fmt.Sprintf("0.0.0.0:%d", s.Port())
	fmt.Println("Runing server on port:", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
