package CoapPubsub

import (
	"log"
	"net"

	"github.com/dustin/go-coap"
)

type chanMapStringList map[*net.UDPAddr][]string
type stringMapChanList map[string][]*net.UDPAddr

type CoapPubsubServer struct {
	capacity int

	msgIndex uint16 //for increase and sync message ID

	//map to store "chan -> Topic List" for find subscription
	clientMapTopics chanMapStringList
	//map to store "topic -> chan List" for publish
	topicMapClients stringMapChanList
}

func NewCoapPubsubServer(maxChannel int) *CoapPubsubServer {
	cSev := new(CoapPubsubServer)
	cSev.capacity = maxChannel
	cSev.clientMapTopics = make(map[*net.UDPAddr][]string, maxChannel)
	cSev.topicMapClients = make(map[string][]*net.UDPAddr, maxChannel)
	cSev.msgIndex = 0
	return cSev
}

func (c *CoapPubsubServer) genMsgID() uint16 {

	c.msgIndex = c.msgIndex + 1
	return c.msgIndex
}

func (c *CoapPubsubServer) subscribe(topic string, client *net.UDPAddr) {
	topicFound := false
	if val, exist := c.topicMapClients[topic]; exist {
		for _, v := range val {
			if v == client {
				topicFound = true
			}
		}
	}
	if topicFound == false {
		c.topicMapClients[topic] = append(c.topicMapClients[topic], client)
	}

	clientFound := false
	if val, exist := c.clientMapTopics[client]; exist {
		for _, v := range val {
			if v == topic {
				clientFound = true
			}
		}
	}

	if clientFound == false {
		c.clientMapTopics[client] = append(c.clientMapTopics[client], topic)
	}
}

func (c *CoapPubsubServer) publish(l *net.UDPConn, topic string, msg string) {
	if clients, exist := c.topicMapClients[topic]; !exist {
		return
	} else { //topic exist, publish it
		for _, client := range clients {
			c.publishMsg(l, client, topic, msg)
		}
	}

}

func (c *CoapPubsubServer) handleCoAPMessage(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	topic := m.Path()[0]
	etag := parseToString(m.Option(coap.ETag))
	log.Printf("Got message path=%q: %#v from %v", m.Path(), m, a)
	log.Println("code=", m.Code, " option=", etag)

	if etag == "SUB" {
		c.subscribe(topic, a)
	} else if etag == "PUB" {
		c.publish(l, topic, string(m.Payload))
	}

	return nil
}

func (c *CoapPubsubServer) ListenAndServe(udpPort string) {
	log.Fatal(coap.ListenAndServe("udp", udpPort,
		coap.FuncHandler(func(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
			return c.handleCoAPMessage(l, a, m)
		})))
}

func (c *CoapPubsubServer) publishMsg(l *net.UDPConn, a *net.UDPAddr, topic string, msg string) {
	m := coap.Message{
		Type:      coap.NonConfirmable,
		Code:      coap.Content,
		MessageID: c.genMsgID(),
		Payload:   []byte(msg),
	}

	m.SetOption(coap.ContentFormat, coap.TextPlain)
	m.SetOption(coap.LocationPath, topic)

	log.Printf("Transmitting %v msg=%s", m, msg)
	err := coap.Transmit(l, a, m)
	if err != nil {
		log.Printf("Error on transmitter, stopping: %v", err)
		return
	}
}
