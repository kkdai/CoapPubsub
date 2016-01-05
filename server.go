package CoapPubsub

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/dustin/go-coap"
)

type chanMapStringList map[chan interface{}][]string
type stringMapChanList map[string][]chan interface{}

type CoapPubsubServer struct {
	capacity int

	//map to store "chan -> Topic List" for find subscription
	clientMapTopics chanMapStringList
	//map to store "topic -> chan List" for publish
	topicMapClients stringMapChanList
}

func NewCoapPubsubServer(maxChannel int) *CoapPubsubServer {
	cSev := new(CoapPubsubServer)
	cSev.capacity = maxChannel
	cSev.clientMapTopics = make(map[chan interface{}][]string, maxChannel)
	cSev.topicMapClients = make(map[string][]chan interface{}, maxChannel)
	return cSev
}

func (c *CoapPubsubServer) handleCoAPMessage(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {

	//Only handle ack msg for pubsub
	if m.Type != coap.Acknowledgement {
		return nil
	}

	log.Printf("Got message path=%q: %#v from %v", m.Path(), m, a)
	log.Println("code=", m.Code, " option=", m.Option(coap.Observe))
	if m.Code == coap.GET && m.Option(coap.Observe) != nil {
		go periodicTransmitter(l, a, m)
	}
	return nil

}

func (c *CoapPubsubServer) ListenAndServe(udpPort string) {
	log.Fatal(coap.ListenAndServe("udp", udpPort,
		coap.FuncHandler(func(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
			return c.handleCoAPMessage(l, a, m)
		})))
}

func periodicTransmitter(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) {
	subded := time.Now()

	for {
		msg := coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: m.MessageID,
			Payload:   []byte(fmt.Sprintf("Been running for %v", time.Since(subded))),
		}

		msg.SetOption(coap.ContentFormat, coap.TextPlain)
		msg.SetOption(coap.LocationPath, m.Path())

		log.Printf("Transmitting %v", msg)
		err := coap.Transmit(l, a, msg)
		if err != nil {
			log.Printf("Error on transmitter, stopping: %v", err)
			return
		}

		time.Sleep(time.Second)
	}
}
