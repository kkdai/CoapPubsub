package CoapPubsub

import (
	"errors"
	"log"
	"time"

	"github.com/dustin/go-coap"
)

type subConnection struct {
	channel   chan string
	clientCon *coap.Conn
}

type CoapPubsubClient struct {
	msgIndex uint16
	serAddr  string
	subList  map[string]subConnection
}

// Create a pubsub client for CoAP protocol
// It will connect to server and make sure it alive and start heart beat
// To keep udp port open, we will send heart beat event to server every minutes
func NewCoapPubsubClient(servAddr string) *CoapPubsubClient {
	c := new(CoapPubsubClient)
	c.subList = make(map[string]subConnection, 0)
	c.serAddr = servAddr

	//TODO: connection check if any error

	//Start heart beat
	go c.heartBeat()
	return c
}

//Add Subscribetion on topic and return a channel for user to wait data
func (c *CoapPubsubClient) AddSub(topic string) (chan string, error) {
	if val, exist := c.subList[topic]; exist {
		//if topic already sub, return and not send to server
		return val.channel, nil
	}

	addSubReq := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: c.getMsgID(),
		Payload:   []byte(""),
	}

	addSubReq.SetOption(coap.ETag, "ADDSUB")
	addSubReq.SetPathString(topic)

	conn, err := coap.Dial("udp", c.serAddr)
	if err != nil {
		log.Printf("ADDSUB>>Error dialing: %v \n", err)
		return nil, errors.New("Dial failed")
	}

	subChan := make(chan string)

	//Send out to server
	conn.Send(addSubReq)
	go c.waitSubResponse(conn, subChan, topic)

	//Add client connection into member variable for heart beat
	clientConn := subConnection{channel: subChan, clientCon: conn}
	c.subList[topic] = clientConn

	return subChan, nil
}

func (c *CoapPubsubClient) waitSubResponse(conn *coap.Conn, ch chan string, topic string) {
	var rv *coap.Message
	var err error
	var keepLoop bool
	keepLoop = true
	for keepLoop {
		if rv != nil {
			if err != nil {
				log.Fatalf("Error receiving: %v", err)
			}
			log.Printf("Got %s", rv.Payload)
		}
		rv, err = conn.Receive()
		log.Println("receiv:", rv, " err=", err)

		if err == nil {
			ch <- string(rv.Payload)
		}

		time.Sleep(time.Second)
		if _, exist := c.subList[topic]; !exist {
			//sub topic already remove, leave loop
			log.Println("Loop topic:", topic, " already remove leave loop")
			keepLoop = false
		}
	}
}

func (c *CoapPubsubClient) RemoveSub(topic string) {
}

func (c *CoapPubsubClient) getMsgID() uint16 {
	c.msgIndex = c.msgIndex + 1
	return c.msgIndex
}

func (c *CoapPubsubClient) heartBeat() {
	log.Println("Starting heart beat loop call")
	hbReq := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: c.getMsgID(),
		Payload:   []byte("Heart beat msg."),
	}

	hbReq.SetOption(coap.ETag, "HB")

	for {

		for k, conn := range c.subList {
			conn.clientCon.Send(hbReq)
			log.Println("Send the heart beat in topic ", k)
		}

		time.Sleep(time.Minute)
	}
}
