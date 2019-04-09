package commands

import (
	"context"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/indrenicloud/tricloud-server/core"
	"github.com/kr/pty"
)

// Terminal implements the command func signature and all terminal related
// logic starts from here
func Terminal(msg *core.MessageFormat, out chan []byte) {

	//terinalLock.Lock()
	log.Println("got a lock")
	//defer terinalLock.Unlock()
	//defer log.Println("freed lock")

	term, ok := terminals[msg.ReceiverConnid]

	if !ok {
		term = newTerminal(msg.ReceiverConnid, out)
		terminals[msg.ReceiverConnid] = term
		term.run()
		return
	}

	data, ok := msg.Arguments["data"]
	log.Println("almost sending to terminal ")
	if ok {
		log.Println("almost sending to terminal ")
		term.inData <- []byte(data)
	}

}

func unregisterTerminal(id core.UID) {
	//pass TODO
}

var terinalLock *sync.Mutex
var terminals map[core.UID]*terminal

func init() {
	terinalLock = &sync.Mutex{}
	terminals = make(map[core.UID]*terminal)
}

type windowinfo struct {
	// title string
	rows int
	col  int
	x    int
	y    int
}

type terminal struct {
	ownerConnID core.UID
	out         chan []byte
	inData      chan []byte
	resize      chan interface{}
	ctx         context.Context
	ctxFunc     context.CancelFunc
	cmd         *exec.Cmd
	tty         *os.File
}

func newTerminal(uid core.UID, outchannel chan []byte) *terminal {
	//c := exec.Command("/bin/bash", "-l")
	//c.Env = append(os.Environ(), "TERM=xterm")
	c := exec.Command("bash")
	tty, err := pty.Start(c)
	if err != nil {
		log.Println("couldnot create terminal", err)
		return nil
	}

	ctx := context.Background()
	ctx, ctxFunc := context.WithCancel(ctx)

	return &terminal{
		ownerConnID: uid,
		out:         outchannel,
		inData:      make(chan []byte),
		resize:      make(chan interface{}),
		ctx:         ctx,
		ctxFunc:     ctxFunc,
		cmd:         c,
		tty:         tty,
	}
}

func (t *terminal) run() {

	//cleanup on exit
	defer func() {
		t.cmd.Process.Kill()
		t.cmd.Process.Wait()
		t.tty.Close()
		unregisterTerminal(t.ownerConnID)
	}()

	rCtx, _ := context.WithCancel(t.ctx)

	//read from terminal and send to server
	go func() {
		for {
			buf := make([]byte, 1024)
			read, err := t.tty.Read(buf)
			if err != nil {
				// notify error may be
				log.Println("error reading from terminal:", err)
				return
			}
			//construct msg with buf[:read]
			log.Println("sending bytes")
			t.out <- ConstructMessage(t.ownerConnID, core.CMD_TERMINAL, []string{string(buf[:read])})

			select {
			case _ = <-rCtx.Done():
				return
			default:
			}
		}
	}()

	for {

		select {
		case _ = <-t.ctx.Done():
			return
		case data := <-t.inData:
			t.tty.Write(data)
		case _ = <-t.resize:
			//pass
		}

	}

}
