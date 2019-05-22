package screenshot

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"

	"../../controller"
	"../../lexer"
	"github.com/kbinani/screenshot"
)

type Screenshoter struct {
	Resources []image.Image
	Control   *controller.Controller
	Path      string
}

type NetWriter struct {
	Control *controller.Controller
}

func (nw *NetWriter) Write(p []byte) (int, error) {
	nw.Control.Send(string(p))
	return len(p), nil
}

func New(control *controller.Controller) *Screenshoter {
	a := &Screenshoter{nil, control, ""}
	control.Interface("screenshot", a)
	return a
}

func TakeScreenshot() ([]image.Image, error) {
	n := screenshot.NumActiveDisplays()
	imgs := make([]image.Image, 0)
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			return nil, err
		}
		imgs = append(imgs, img)
	}
	return imgs, nil
}

func (sc *Screenshoter) Send(tk []lexer.Token) {
	fmt.Println(lexer.MakeString(tk))
	if len(tk) > 0 {
		switch tk[0].Value {
		case "take":
			lexer.AllowLength(tk, 1, func() error {
				res, e := TakeScreenshot()
				if e == nil {
					sc.Resources = res
				}
				return e
			})
		case "send":
			lexer.AllowLength(tk, 1, func() error {
				sc.Encode(true)
				return nil
			})
		case "path":
			lexer.AllowLength(tk, 2, func() error {
				sc.Path = tk[1].Value
				return nil
			})
		case "save":
			lexer.AllowAtLeastLength(tk, 1, func() error {
				if len(tk) > 1 {
					sc.Path = tk[1].Value
				}
				sc.Encode(false)
				return nil
			})
		}
	}
}

func (sc *Screenshoter) Encode(send bool) {
	var opt jpeg.Options
	opt.Quality = 100
	var wr io.Writer
	var err error
	if send {
		wr = &NetWriter{}
	} else {
		wr, err = os.Create(sc.Path)
	}
	if err == nil {
		for _, res := range sc.Resources {
			if res != nil {
				jpeg.Encode(wr, res, &opt)
				sc.Control.Send("\n")
			}
		}
	}
}

func (sc *Screenshoter) Start() {

}
