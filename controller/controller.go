package controller

import (
	"fmt"
	"net"

	"../lexer"
	"../remotecontrol"
)

type Include interface {
	Send([]lexer.Token)
	Start()
}

type Controller struct {
	Network         *remotecontrol.Remote
	Includes        map[string]Include
	ActiveInterface string
	port            int
	address         string
	DataController  func([]lexer.Token)
}

func New() *Controller {
	cnt := Controller{
		nil,
		make(map[string]Include),
		"",
		-1,
		"",
		func(data []lexer.Token) {},
	}
	cnt.Network = remotecontrol.New(&cnt)
	cnt.DataController = func(data []lexer.Token) {
		switch data[0].Value {
		case "interface":
			lexer.AllowLength(data, 2, func() error {
				cnt.ActiveInterface = data[1].Value
				return nil
			})
		default:
			if len(cnt.ActiveInterface) > 0 {
				if val, ok := cnt.Includes[cnt.ActiveInterface]; ok {
					val.Send(data)
				}
			}
		}
	}
	return &cnt
}

func (control *Controller) Port(port int) {
	control.port = port
}

func (control *Controller) Address(address string) {
	control.address = address
}

func (control *Controller) Interface(name string, object Include) {
	control.Includes[name] = object
}

func (control *Controller) Send(content string) {
	if control.Network.InnerClient.Can() {
		control.Network.InnerClient.Send(content)
	}
}

func (control *Controller) Read(content string) {
	a, e := lexer.ParseLine(content)
	if e == nil {
		control.DataController(a)
	}
}

func (control *Controller) Connected(c net.Conn) {

}

func (control *Controller) Start() {
	if control.port > 0 {
		control.Network.Listen(control.port, control.DataController)
	}
	if len(control.address) > 0 {
		e := control.Network.Connect(control.address)
		if e != nil {
			fmt.Println(e.Error())
		}
		control.Network.Get(control.DataController)
	}
}
