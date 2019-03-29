package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// Connection encapsulates websocket conn and related stuff
type Connection struct {
	conn *websocket.Conn
	Out  chan []byte
	In   chan []byte

	// reader/writer notifies main of error through ErrorChannel
	ErrorChannel chan struct{}
	// reader/writer contex
	readerctx context.Context
	writerctx context.Context
}

var addr = flag.String("addr", "localhost:8081", "http service address")

const (
	// Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * 9) / 10

	// Maximum message size allowed from peer.
	MaxMessageSize = 512
)

// NewConnection is constructor
func NewConnection(ctx context.Context, In, Out chan []byte, ErrorChannel chan struct{}) *Connection {

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return nil
	}

	newctx1, _ := context.WithCancel(ctx)
	newctx2, _ := context.WithCancel(ctx)

	return &Connection{
		conn:         c,
		In:           In,
		Out:          Out,
		ErrorChannel: ErrorChannel,
		readerctx:    newctx1,
		writerctx:    newctx2,
	}
}

// Reader reads message/command from conn and gives to worker
func (c *Connection) Reader() {

	for {

		_, message, err := c.conn.ReadMessage()

		log.Println("reading connection")

		if err != nil {
			log.Println("read:", err)
			c.ErrorChannel <- struct{}{}
			return
		}
		log.Printf("recv: %s", message)

		c.In <- message //sending to worker coroutine

		select {
		// checking if someone want to close reader
		case _ = <-c.readerctx.Done():
			return
		default:
		}

	}
}

func (c *Connection) Writer() {

	ticker := time.NewTicker(PongWait)

	defer ticker.Stop()

	c.conn.SetPongHandler(func(appData string) error {
		ticker.Stop()
		ticker = time.NewTicker(PongWait)

		return nil
	})

	for {

		select {
		// checking if someone want to close writer
		case _ = <-c.writerctx.Done():
			return
		case sendData := <-c.Out:
			log.Println("writing connection")
			err := c.conn.WriteMessage(websocket.TextMessage, []byte(sendData))
			if err != nil {
				log.Println("write:", err)
				c.ErrorChannel <- struct{}{}
				return

			}
		case _ = <-ticker.C:
			log.Println("Pinging server")
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println("write:", err)
				c.ErrorChannel <- struct{}{}
				return

			}
		}

	}

}

func (c *Connection) Run() {
	go c.Reader()
	go c.Writer()
}
