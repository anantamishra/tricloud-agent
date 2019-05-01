package conn

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	readerctx    context.Context
	readctxFunc  context.CancelFunc
	writerctx    context.Context
	writectxFunc context.CancelFunc
}

var addr = flag.String("addr", "localhost:8081", "http service address")

func NewConnection(ctx context.Context, In, Out chan []byte, ErrorChannel chan struct{}) *Connection {

	cf := GetConfig()

	if cf.UUID == "" {
		if cf.ApiKey == "" {
			panic("Need api key")
		}

		client := &http.Client{}
		req, _ := http.NewRequest("POST", "http://localhost:8081/registeragent", nil)
		req.Header.Add("Api-key", cf.ApiKey)

		resp, err := client.Do(req)

		if err != nil {
			panic("server error")
		}
		body, err := ioutil.ReadAll(resp.Body)

		resbody := map[string]string{}
		json.Unmarshal(body, &resbody)

		uid := resbody["data"]
		logg.Log("My ID:", uid)

		if uid == "" {
			panic("Server didnot register us, every man for himself")
		}

		cf.UUID = uid
		SaveConfig(cf)
	}

	updateSystemInfo(cf.UUID)

	u := url.URL{Scheme: "ws", Host: *addr, Path: fmt.Sprintf("/websocket/%s", cf.UUID)}
	logg.Log("connecting to :", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logg.Log("ERROR", "dial", err)
		return nil
	}

	newctx1, ctx1func := context.WithCancel(ctx)
	newctx2, ctx2func := context.WithCancel(ctx)

	return &Connection{
		conn:         c,
		In:           In,
		Out:          Out,
		ErrorChannel: ErrorChannel,
		readerctx:    newctx1,
		readctxFunc:  ctx1func,
		writerctx:    newctx2,
		writectxFunc: ctx2func,
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

	defer c.conn.Close()

	for {

		select {
		case _ = <-c.writerctx.Done(): // checking if someone want to close writer
			return
		case sendData := <-c.Out:

			logg.Log("Writing to connection")

			err := c.conn.WriteMessage(websocket.BinaryMessage, []byte(sendData))

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
