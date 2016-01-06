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

func (c *CoapPubsubServer) removeSubscribe(topic string, client *net.UDPAddr) {
	removeIndexT2C := -1
	if val, exist := c.topicMapClients[topic]; exist {
		for k, v := range val {
			if v == client {
				removeIndexT2C = k
			}
		}
		if removeIndexT2C != -1 {
			sliceClients := c.topicMapClients[topic]
			if len(sliceClients) > 1 {
				c.topicMapClients[topic] = append(sliceClients[:removeIndexT2C], sliceClients[removeIndexT2C+1:]...)
			} else {
				delete(c.topicMapClients, topic)
			}
		}
	}

	removeIndexC2T := -1
	if val, exist := c.clientMapTopics[client]; exist {
		for k, v := range val {
			if v == topic {
				removeIndexC2T = k
			}
		}
		if removeIndexC2T != -1 {
			sliceTopics := c.clientMapTopics[client]
			if len(sliceTopics) > 1 {
				c.clientMapTopics[client] = append(sliceTopics[:removeIndexC2T], sliceTopics[removeIndexC2T+1:]...)
			} else {
				delete(c.clientMapTopics, client)
			}
		}
	}

}

func (c *CoapPubsubServer) addSubscribe(topic string, client *net.UDPAddr) {
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
			log.Println("topic->", topic, " PUB to ", client, " msg=", msg)
		}
	}
	log.Println("pub finished")
}

func (c *CoapPubsubServer) handleCoAPMessage(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	topic := m.Path()[0]
	cmd := parseToString(m.Option(coap.ETag))
	log.Println("cmd=", cmd, " topic=", topic, " msg=", string(m.Payload))
	log.Println("code=", m.Code, " option=", cmd)

	if cmd == "SUB" {
		log.Println("sub topic=", topic, " in client=", a)
		c.addSubscribe(topic, a)
		c.responseOK(l, a, m)

	} else if cmd == "PUB" {
		c.publish(l, topic, string(m.Payload))
		c.responseOK(l, a, m)
	}

	for k, v := range c.topicMapClients {
		log.Println("Topic=", k, " sub by client=>", v)
	}
	return nil
}

func (c *CoapPubsubServer) ListenAndServe(udpPort string) {
	log.Fatal(coap.ListenAndServe("udp", udpPort,
		coap.FuncHandler(func(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
			return c.handleCoAPMessage(l, a, m)
		})))
}

func (c *CoapPubsubServer) responseOK(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) {
	m2 := coap.Message{
		Type:      coap.Acknowledgement,
		Code:      coap.Content,
		MessageID: m.MessageID,
		Payload:   m.Payload,
	}

	m2.SetOption(coap.ContentFormat, coap.TextPlain)
	m2.SetOption(coap.LocationPath, m.Path())

	err := coap.Transmit(l, a, m2)
	if err != nil {
		log.Printf("Error on transmitter, stopping: %v", err)
		return
	}
}

func (c *CoapPubsubServer) publishMsg(l *net.UDPConn, a *net.UDPAddr, topic string, msg string) {
	m := coap.Message{
		Type:      coap.Confirmable,
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
