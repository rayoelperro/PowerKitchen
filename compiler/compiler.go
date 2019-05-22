package compiler

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	. "../lexer"
	. "../libraries"
)

var Commands = []string{"include", "set", "do", "insert", "import", "address", "lansearch", "port", "visible"}

type ResultFile struct {
	Head       string
	Body       string
	Controller string
	Visible    bool
}

func NewResultFile() *ResultFile {
	var head = `package main

import (
%s
)`
	var body = `func main() {
%s
	for true {
		continue
	}
}`
	rf := ResultFile{head, body, "control", false}
	rf.Line(rf.Controller + " := controller.New()")
	return &rf
}

func (f *ResultFile) Library(imp string) {
	f.Head = fmt.Sprintf(f.Head, "	\""+imp+"\"\n%s")
}

func (f *ResultFile) Line(ln string) {
	f.Body = fmt.Sprintf(f.Body, "	"+ln+"\n%s")
}

func (f *ResultFile) Output() string {
	return fmt.Sprintf(f.Head+"\n\n"+f.Body, "	\"../controller\"", "	"+f.Controller+".Start()")
}

func Compile(file *PWKTFile, cwd string, exd string) error {
	rsfl := NewResultFile()
	libs, e := Libs(filepath.Join(exd, "libs.json"))
	if e != nil {
		return e
	}
	for _, ln := range file.Lines {
		e := UseLine(rsfl, ln, libs)
		if e != nil {
			return e
		}
	}
	path := filepath.Join(exd, "build", file.Name()+".go")
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		e := os.Mkdir(filepath.Dir(path), os.ModePerm)
		if e != nil {
			return e
		}
	}
	ioutil.WriteFile(path, []byte(rsfl.Output()), os.ModePerm)
	er := Executable(path, rsfl.Visible)
	time.Sleep(3 * time.Second)
	if er == nil {
		expath := filepath.Join(exd, "build", file.Name()+".exe")
		nepath := filepath.Join(cwd, file.Name()+".exe")
		e := os.Rename(expath, nepath)
		if e == nil {
			Info.Println("File saved at: " + nepath)
		} else {
			Error.Println(e.Error())
			Info.Println("File saved at: " + expath)
		}
	}
	os.Remove(path)
	return er
}

func UseLine(rsfl *ResultFile, ln []Token, libs *LibCollection) error {
	switch ln[0].Value {
	case "include":
		return AllowLength(ln, 2, func() error {
			lb, e := libs.Take(ln[1].Value)
			if e != nil {
				return e
			}
			if lb.Standard {
				rsfl.Library(lb.Load)
				rsfl.Line(lb.Var + " := " +
					filepath.Base(lb.Load) + ".New(" + rsfl.Controller + ")")
				rsfl.Line(lb.Var + ".Start()")
			}
			return nil
		})
	case "set":
		return AllowLength(ln, 4, func() error {
			lb, e := libs.Take(ln[1].Value)
			if e != nil {
				return e
			}
			for _, p := range lb.Properties {
				if p == ln[2].Value {
					rsfl.Line(lb.Var + "." + p + " = " + ln[3].Value)
					return nil
				}
			}
			return fmt.Errorf("The property %s in the library %s does not exists", ln[2].Value, ln[1].Value)
		})
	case "do":
		return AllowLength(ln, 3, func() error {
			lb, e := libs.Take(ln[1].Value)
			if e != nil {
				return e
			}
			for _, m := range lb.Methods {
				if m.Name == ln[2].Value {
					rsfl.Line(m.Code)
					return nil
				}
			}
			return fmt.Errorf("The method %s in the library %s does not exists", ln[2].Value, ln[1].Value)
		})
	case "insert":
		//Not yet
	case "import":
		//Not yet
	case "address":
		return AllowLength(ln, 2, func() error {
			rsfl.Line(rsfl.Controller + ".Address(\"" + ln[1].Value + "\")")
			return nil
		})
	case "lansearch":
		//Not yet
	case "port":
		return AllowLength(ln, 2, func() error {
			rsfl.Line(rsfl.Controller + ".Port(" + ln[1].Value + ")")
			return nil
		})
	case "visible":
		return AllowLength(ln, 2, func() error {
			var v bool
			_, e := fmt.Sscan(ln[1].Value, &v)
			if e == nil {
				rsfl.Visible = v
			}
			return e
		})
	}
	return nil
}

func Executable(gofile string, visible bool) error {
	arguments := []string{"go", "build"}
	if runtime.GOOS == "linux" {
		arguments = append([]string{"GOOS=windows", "GOARCH=386"}, arguments...)
	}
	if !visible {
		arguments = append(arguments, "-ldflags", "-H=windowsgui")
	}
	arguments = append(arguments, filepath.Base(gofile))
	cmd := exec.Command(arguments[0], arguments[1:]...)
	cmd.Dir = filepath.Dir(gofile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}
