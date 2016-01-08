package main

import (
	"log"

	. "github.com/kkdai/CoapPubsub"
)

func main() {
	log.Println("Server start....")
	serv := NewCoapPubsubServer(1024)
	serv.ListenAndServe(":5683")
}
