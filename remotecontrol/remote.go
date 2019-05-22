package remotecontrol

import (
	"errors"
	"net"
	"strconv"

	. "../lexer"
)

type NetworkReceptor interface {
	Connected(net.Conn)
	Read(string)
}

type Remote struct {
	InnerServer     *Server
	InnerClient     *Client
	NetworkReceiver NetworkReceptor
}

func New(n NetworkReceptor) *Remote {
	return &Remote{
		&Server{},
		&Client{},
		n,
	}
}

var Commands = []string{"connect", "end", "interface", "send", "open"}

func (rmt *Remote) Command(tk []Token) error {
	switch tk[0].Value {
	case "connect":
		return AllowLength(tk, 2, func() error {
			if rmt.NetworkReceiver != nil {
				rmt.InnerClient.Get(func(s string) error {
					rmt.NetworkReceiver.Read(s)
					return nil
				})
			}
			return rmt.Connect(tk[1].Value)
		})
	case "end":
		return AllowLength(tk, 1, func() error {
			return rmt.CloseConnection()
		})
	case "interface":
		return AllowLength(tk, 2, func() error {
			rmt.Send("interface " + tk[1].Value)
			return nil
		})
	case "send":
		return AllowAtLeastLength(tk, 2, func() error {
			rmt.Send(MakeString(tk[1:]))
			return nil
		})
	case "open":
		return AllowLength(tk, 2, func() error {
			if rmt.NetworkReceiver != nil {
				v, e := strconv.Atoi(tk[1].Value)
				if e != nil {
					return e
				}
				return rmt.InnerServer.AdvanceStart(v, func(s string) error {
					rmt.NetworkReceiver.Read(s)
					return nil
				}, rmt.NetworkReceiver.Connected)
			} else {
				return errors.New("You must create an emisor for open port")
			}
		})
	}
	return nil
}

func (rmt *Remote) Connect(address string) error {
	return rmt.InnerClient.Connect(address)
}

func (rmt *Remote) CloseConnection() error {
	if rmt.InnerClient.Can() {
		return rmt.InnerClient.Close()
	}
	return nil
}

func (rmt *Remote) Send(message string) {
	if rmt.InnerClient.Can() {
		rmt.InnerClient.Send(message)
	}
	rmt.InnerServer.SendAll(message)
}

func (rmt *Remote) Get(fn func([]Token)) error {
	return rmt.InnerClient.Get(func(msg string) error {
		m, e := ParseLine(msg)
		if e != nil {
			return e
		}
		fn(m)
		return nil
	})
}

func (rmt *Remote) Listen(port int, fn func([]Token)) error {
	return rmt.InnerServer.Start(port, func(msg string) error {
		m, e := ParseLine(msg)
		if e != nil {
			return e
		}
		fn(m)
		return nil
	})
}
