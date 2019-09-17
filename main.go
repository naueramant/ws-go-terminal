package main

import (
	"io"
	"net/http"
	"os/exec"

	"github.com/kr/pty"
	"golang.org/x/net/websocket"
)

// TTYServer handles the websocket tty connections
func TTYServer(ws *websocket.Conn) {
	cmd := exec.Command("/bin/bash")
	tty, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}

	defer tty.Close()

	// pipe to/from websocket to the TTY:
	go func() {
		io.Copy(ws, tty)
	}()
	go func() {
		io.Copy(ws, tty)
	}()
	go func() {
		io.Copy(tty, ws)
	}()

	// wait for the command to exit, then close the websocket
	cmd.Wait()
}

func main() {
	http.Handle("/shell", websocket.Handler(TTYServer))
	http.Handle("/", http.FileServer(http.Dir(".")))

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
