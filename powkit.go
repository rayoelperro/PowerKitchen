package main

import (
	"flag"
	"os"
	"path/filepath"

	"./compiler"

	"github.com/lukesampson/figlet/figletlib"
)

func main() {
	var source = flag.String("s", "", "It will compile the PowerKitchen sourcefile (.pwkt) into a .exe")
	var terminal = flag.Bool("c", false, "If it is true, the console will show the PowerKitchen command console")
	var module = flag.String("m", "", "Follows the instructions of a .json file to create a module (.mbpk) file")
	var output = flag.String("o", "", "It will set the output file")
	flag.Parse()
	cwd, _ := os.Getwd()
	exd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := filepath.Join(exd, "fonts")
	f, err := figletlib.GetFontByName(path, "standard")
	if err != nil {
		f = nil
	} else if *terminal {
		figletlib.PrintMsg("PowerKitchen", f, 80, f.Settings(), "center")
	}
	if len(*source) > 0 {
		file, e := compiler.OpenPWKTFile(*source)
		if e == nil {
			e := compiler.Compile(file, cwd, exd)
			if e != nil {
				compiler.Error.Println(e.Error())
			}
		} else {
			compiler.Error.Println(e.Error())
		}
	}
	if len(*module) > 0 {

	}
	if *terminal {
		if len(*output) < 1 {
			*output = "cooked"
		}
		compiler.StartConsole(cwd, exd, *output)
	}
}
