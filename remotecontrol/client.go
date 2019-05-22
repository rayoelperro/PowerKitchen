package remotecontrol

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	Connection net.Conn
	Host       string
}

func NewClient() *Client {
	return &Client{
		Connection: nil,
		Host:       "",
	}
}

func (client *Client) Connect(ip string) error {
	c, e := net.Dial("tcp", ip)
	if e != nil {
		return e
	}
	client.Connection = c
	client.Host = ip
	return nil
}

func (client *Client) Send(data string) {
	fmt.Fprint(client.Connection, data+"\n")
}

func (client *Client) Close() error {
	e := client.Connection.Close()
	client.Connection = nil
	return e
}

func (client *Client) Get(fn func(string) error) error {
	go func() {
		buf := bufio.NewReader(client.Connection)
		for {
			msg, e := buf.ReadString('\n')
			if e != nil {
			back:
				if client.Connect(client.Host) != nil {
					goto back
				} else {
					buf = bufio.NewReader(client.Connection)
				}
			} else {
				fn(strings.Trim(msg, "\n"))
			}
		}
	}()
	return nil
}

func (client *Client) Can() bool {
	return client.Connection != nil
}
