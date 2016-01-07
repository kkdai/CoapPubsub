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
	serAddr string
	subList map[string]subConnection
}

func NewCoapPubsubClient(servAddr string) *CoapPubsubClient {
	c := new(CoapPubsubClient)
	c.subList = make(map[string]subConnection, 0)
	c.serAddr = servAddr

	//connection check if any error

	//Start heart beat
	go c.heartBeat()
	return c
}

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
	go c.waitSubResponse(conn, subChan)

	//Add client connection into member variable for heart beat
	clientConn := subConnection{channel: subChan, clientCon: conn}
	c.subList[topic] = clientConn

	return subChan, nil
}

func (c *CoapPubsubClient) waitSubResponse(conn *coap.Conn, ch chan string) {
	var rv *coap.Message
	var err error
	for {
		if rv != nil {
			if err != nil {
				log.Fatalf("Error receiving: %v", err)
			}
			log.Printf("Got %s", rv.Payload)
		}
		rv, err = conn.Receive()
		log.Println("receiv:", rv, " err=", err)
		ch <- string(rv.Payload)

		time.Sleep(time.Second)
	}
}

func (c *CoapPubsubClient) RemoveSub(topic string) {
}

func (c *CoapPubsubClient) getMsgID() uint16 {
	return 12345
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
