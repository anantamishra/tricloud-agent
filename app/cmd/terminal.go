package cmd

import (
	"context"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/indrenicloud/tricloud-agent/app/logg"
	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/kr/pty"
)

// Terminal implements the command func signature and all terminal related
// logic starts from here
func Terminal(rawdata []byte, out chan []byte, ctx context.Context) {

	//terinalLock.Lock()
	logg.Log("got a lock")
	//defer terinalLock.Unlock()
	//defer log.Println("freed lock")
	termdata := &wire.TermData{}

	header, err := wire.Decode(rawdata, termdata)
	if err != nil {
		return
	}

	term, ok := terminals[header.Connid]

	if !ok {
		term = newTerminal(header.Connid, out)
		terminals[header.Connid] = term
		term.run()
		return
	}

	logg.Log("Almost sending to terminal")
	term.inData <- []byte(termdata.Data)
	logg.Log("Sent to terminal ")
	return

}

func unregisterTerminal(id wire.UID) {

	// todo lock
	term, ok := terminals[id]
	if !ok {
		logg.Log("terminal already cancelled")
		return
	}
	if term.ctx.Err() == nil {
		term.ctxFunc()
	}
	delete(terminals, id)

}

var terinalLock *sync.Mutex
var terminals map[wire.UID]*terminal

func init() {
	terinalLock = &sync.Mutex{}
	terminals = make(map[wire.UID]*terminal)
}

type windowinfo struct {
	// title string
	rows int
	col  int
	x    int
	y    int
}

type terminal struct {
	ownerConnID wire.UID
	out         chan []byte
	inData      chan []byte
	resize      chan interface{}
	ctx         context.Context
	ctxFunc     context.CancelFunc
	cmd         *exec.Cmd
	tty         *os.File
}

func newTerminal(uid wire.UID, outchannel chan []byte) *terminal {
	//c := exec.Command("/bin/bash", "-l")
	//c.Env = append(os.Environ(), "TERM=xterm")
	c := exec.Command("bash")
	tty, err := pty.Start(c)
	if err != nil {
		logg.Log("couldnot create terminal", err)
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
		unregisterTerminal(t.ownerConnID)
		t.cmd.Process.Kill()
		t.cmd.Process.Wait()
		t.tty.Close()
	}()

	rCtx, _ := context.WithCancel(t.ctx)

	//read from terminal and send to server
	go func() {
		for {
			buf := make([]byte, 1024)
			read, err := t.tty.Read(buf)
			if err != nil {
				// notify error may be
				logg.Log("error reading from terminal:", err)
				return
			}

			logg.Log("sending bytes")

			h := wire.NewHeader(t.ownerConnID, wire.CMD_TERMINAL, wire.AgentToUser)
			outbyte := wire.AttachHeader(h, buf[:read])

			t.out <- outbyte

			select {
			case _ = <-rCtx.Done():
				return
			default:
			}
		}
	}()

	timer := time.NewTimer(time.Minute * 1)

	for {

		select {
		case _ = <-t.ctx.Done():
			return
		case data := <-t.inData:
			timer.Reset(time.Minute * 1)
			t.tty.Write(data)
		case _ = <-t.resize:
			//pass
		case _ = <-timer.C:
			return
		}

	}

}
