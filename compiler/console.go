package compiler

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	. "../lexer"
	"../libraries"
	"../remotecontrol"

	"github.com/fatih/color"
)

type Console struct {
	Current          *PWKTFile
	Network          *remotecontrol.Remote
	System           string
	WorkingDirectory string
	AppDirectory     string
}

var Error *color.Color = color.New(color.FgRed, color.Bold)
var Info *color.Color = color.New(color.FgGreen, color.Bold)
var NetInfo *color.Color = color.New(color.FgYellow, color.Bold)
var LineInp *color.Color = color.New(color.FgBlue, color.Bold)

func StartConsole(workingdir string, appdir string, filename string) {
	con := Console{&PWKTFile{make([][]Token, 0), filename}, nil, runtime.GOOS, workingdir, appdir}
	con.Network = remotecontrol.New(&con)
	scanner := bufio.NewScanner(os.Stdin)
	LineInp.Print(">>>")
	fmt.Print(" ")
	for scanner.Scan() {
		line := scanner.Text()
		p, e := ParseLine(line)
		if e == nil {
			e := Run(p, &con)
			if e != nil {
				Error.Println(e.Error())
			}
		} else {
			Error.Println(e.Error())
		}
		LineInp.Print(">>>")
		fmt.Print(" ")
	}
}

func (Console *Console) Read(msg string) {
	NetInfo.Println(msg)
}

func (Console *Console) Connected(c net.Conn) {
	NetInfo.Println("New connection from: " + c.RemoteAddr().String())
}

var GlobalCommands = []string{"restart", "build", "comment", "clear", "output", "exit", "inspect", "getall"}

func (console *Console) Command(tk []Token) error {
	switch tk[0].Value {
	case "restart":
		return AllowLength(tk, 1, func() error {
			console.Current.Clear()
			return nil
		})
	case "build":
		return AllowLength(tk, 1, func() error {
			return Compile(console.Current, console.WorkingDirectory, console.AppDirectory)
		})
	case "comment":
		return nil
	case "clear":
		return AllowLength(tk, 1, func() error {
			var cmd *exec.Cmd
			if console.System == "windows" {
				cmd = exec.Command("cmd", "/c", "cls")
			} else if console.System == "linux" {
				cmd = exec.Command("clear")
			} else {
				return errors.New("You can not clear this platform")
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stdout
			cmd.Start()
			cmd.Wait()
			return nil
		})
	case "output":
		return AllowLength(tk, 2, func() error {
			console.Current.Path = tk[1].Value
			return nil
		})
	case "exit":
		return AllowLength(tk, 1, func() error {
			os.Exit(-1)
			return nil
		})
	case "inspect":
		return AllowLength(tk, 2, func() error {
			col, e := libraries.Libs(filepath.Join(console.AppDirectory, "libs.json"))
			if e != nil {
				return e
			}
			l, e := col.Take(tk[1].Value)
			if e != nil {
				return e
			}
			Info.Printf("Library %s:\nNeeds: %v\nMethods: %v\nProperties: %v\nIs Standard: %t\nInformation: %v\n", l.Name, l.Requirements, l.MethodsName(), l.Properties, l.Standard, l.Info)
			return nil
		})
	case "getall":
		return AllowLength(tk, 2, func() error {
			col, e := libraries.Libs(filepath.Join(console.AppDirectory, "libs.json"))
			if e != nil {
				return e
			}
			l, e := col.Take(tk[1].Value)
			if e != nil {
				return e
			}
			for _, req := range l.Requirements {
				var cmd *exec.Cmd
				if console.System == "windows" {
					cmd = exec.Command("cmd", "/c", req)
				} else if console.System == "linux" {
					cmd = exec.Command("sh", "-c", req)
				} else {
					return errors.New("You can not call that on this platform")
				}
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stdout
				cmd.Start()
				cmd.Wait()
			}
			return nil
		})
	}
	return nil
}
