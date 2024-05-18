package main

import (
	"flag"
)

func main() {
	port := flag.Uint("port", 5000, "TCP Port Number for BC Server")
	flag.Parse()

	app := NewServer(uint16(*port))
	app.Run()
}
