package main

import (
	"flag"
	"fmt"
	"log"

	. "github.com/kkdai/CoapPubsub"
)

func main() {
	flag.Parse()
	cmd := flag.Arg(0)
	topic := flag.Arg(1)
	msg := flag.Arg(2)

	fmt.Println(cmd, topic, msg)

	client := NewCoapPubsubClient("localhost:5683")
	if client == nil {
		log.Fatalln("Cannot connect to server, please check your setting.")
	}

	if cmd == "ADDSUB" {
		ch, err := client.AddSub(topic)
		log.Println(" ch:", ch, " err=", err)
		log.Println("Got pub from topic:", topic, " pub:", <-ch)
	}
	log.Println("Done")
}
