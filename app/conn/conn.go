package conn

import (
	"context"
	"flag"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/indrenicloud/tricloud-agent/app/logg"
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

// NewConnection is constructor
func NewConnection(ctx context.Context, In, Out chan []byte, ErrorChannel chan struct{}) *Connection {

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/websocket/456456"}
	logg.Log("connecting to :", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logg.Log("ERROR", "dial", err)
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
func (c *Connection) reader() {

	for {

		_, message, err := c.conn.ReadMessage()

		logg.Log("reading connection")

		if err != nil {
			log.Println("read:", err)
			c.ErrorChannel <- struct{}{}
			return
		}
		logg.Log("recv:", message)

		c.In <- message //sending to worker coroutine

		select {
		case _ = <-c.readerctx.Done(): // checking if someone want to close reader
			return
		default:
		}

	}
}

func (c *Connection) writer() {

	defer c.conn.Close()

	for {

		select {
		case _ = <-c.writerctx.Done(): // checking if someone want to close writer
			return
		case sendData := <-c.Out:

			logg.Log("Writing to connection")

			err := c.conn.WriteMessage(websocket.TextMessage, []byte(sendData))

			if err != nil {
				logg.Log("Write Error:", err)
				c.ErrorChannel <- struct{}{}
				return

			}
		}

	}

}

func (c *Connection) Run() {
	logg.Log("Writing Reader, Writer coroutines")
	go c.reader()
	go c.writer()
}

func (c *Connection) close() {
	logg.Log("Closing connection")
	c.conn.Close()
}
