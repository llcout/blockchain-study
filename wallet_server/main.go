package main

import "flag"

func main() {
	port := flag.Uint("port", 8080, "TCP Port number for Wallet Server")
	gateway := flag.String("gateway", "http://127.0.0.1:5000", "Blockchain gateway")
	flag.Parse()

	app := NewWalletServer(uint16(*port), *gateway)
	app.Run()
}
