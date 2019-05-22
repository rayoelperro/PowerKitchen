package io

import (
	"../../controller"
	"../../lexer"
)

type FileControl struct {
	Control *controller.Controller
	Path    string
}

func New(control *controller.Controller) *FileControl {
	a := &FileControl{control, ""}
	control.Interface("io", a)
	return a
}

func (cmd *FileControl) Send(tk []lexer.Token) {

}

func (cmd *FileControl) Start() {

}
