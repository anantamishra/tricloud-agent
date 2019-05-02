package conn

import (
	"context"
	"fmt"
	"log"

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
	readerctx     context.Context
	readctxFunc   context.CancelFunc
	writerctx     context.Context
	writectxFunc  context.CancelFunc
	workerctx     context.Context
	workerctxFunc context.CancelFunc
}

func NewConnection(ctx context.Context, ErrorChannel chan struct{}) *Connection {
	cf := GetConfig()

	//u := url.URL{Scheme: "ws", Host: cf.Url, Path: fmt.Sprintf("/websocket/%s", cf.UUID)}
	u := fmt.Sprintf("ws://%s/websocket/%s", cf.Url, cf.UUID)
	logg.Log("connecting to :", u)

	c, _, err := websocket.DefaultDialer.Dial(u, nil)

	if err != nil {
		logg.Log("ERROR", "dial", err)
		return nil
	}

	newctx1, ctx1func := context.WithCancel(ctx)
	newctx2, ctx2func := context.WithCancel(ctx)
	newctx3, ctx3func := context.WithCancel(ctx)

	return &Connection{
		conn:          c,
		In:            make(chan []byte),
		Out:           make(chan []byte),
		ErrorChannel:  ErrorChannel,
		readerctx:     newctx1,
		readctxFunc:   ctx1func,
		writerctx:     newctx2,
		writectxFunc:  ctx2func,
		workerctx:     newctx3,
		workerctxFunc: ctx3func,
	}
}

// Reader reads message/command from conn and gives to worker
func (c *Connection) reader() {

	defer c.Close()

	for {

		_, message, err := c.conn.ReadMessage()

		logg.Log("reading connection")

		if err != nil {
			log.Println("read:", err)
			c.ErrorChannel <- struct{}{}
			return
		}
		logg.Log("recv:", string(message))

		c.In <- message //sending to worker coroutine

		select {
		case _ = <-c.readerctx.Done(): // checking if someone want to close reader
			return
		default:
		}

	}
}

func (c *Connection) writer() {

	defer c.Close()

	for {

		select {
		case _ = <-c.writerctx.Done(): // checking if someone want to close writer
			return
		case sendData := <-c.Out:

			logg.Log("Writing to connection")

			err := c.conn.WriteMessage(websocket.BinaryMessage, sendData)

			if err != nil {
				logg.Log("Write Error:", err)
				c.ErrorChannel <- struct{}{}
				return

			}
		}

	}

}

func (c *Connection) Run() {
	logg.Log("Running Reader, Writer coroutines")
	go c.reader()
	go c.writer()
	go c.Worker()
}

func (c *Connection) Close() {
	logg.Log("Closing connection")
	if c.readerctx.Err() == nil {
		c.readctxFunc()
	}

	if c.writerctx.Err() == nil {
		c.writectxFunc()
	}

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
