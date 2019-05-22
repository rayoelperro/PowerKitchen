package remotecontrol

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	Listener    net.Listener
	Connections []net.Conn
}

func NewServer() *Server {
	return &Server{
		Listener:    nil,
		Connections: make([]net.Conn, 0),
	}
}

func (server *Server) AdvanceStart(port int, fn func(string) error, nu func(net.Conn)) error {
	listener, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	if e != nil {
		return e
	}
	server.Listener = listener
	go func() {
		for {
			conn, e := server.Listener.Accept()
			if e != nil {
				continue
			}
			server.Connections = append(server.Connections, conn)
			nu(conn)
			go func() {
				buf := bufio.NewReader(conn)
				for {
					msg, e := buf.ReadString('\n')
					if e != nil {
						continue
					}
					fn(strings.Trim(msg, "\n"))
				}
			}()
		}
	}()
	return nil
}

func (server *Server) Start(port int, fn func(string) error) error {
	return server.AdvanceStart(port, fn, func(net.Conn) {})
}

func (server *Server) SendAll(msg string) {
	for _, c := range server.Connections {
		if c != nil {
			fmt.Fprint(c, msg+"\n")
		}
	}
}
