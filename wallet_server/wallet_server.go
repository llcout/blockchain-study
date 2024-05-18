package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{
		port:    port,
		gateway: gateway,
	}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("/workspaces/blockchain/wallet_server/templates/index.html")
		log.Println(err)
		t.Execute(w, "")
	default:
		log.Printf("ERROR: Invalid Method")
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)

	address := fmt.Sprintf("0.0.0.0:%d", ws.Port())
	fmt.Println("Runing server on port:", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
