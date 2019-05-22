package compiler

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	. "../lexer"
)

type PWKTFile struct {
	Lines [][]Token
	Path  string
}

func (file *PWKTFile) Clear() {
	file.Lines = make([][]Token, 0)
}

func (file *PWKTFile) Enter(ts []Token) error {
	file.Lines = append(file.Lines, ts)
	return nil
}

func (file *PWKTFile) Name() string {
	return filepath.Base(file.Path[0 : len(file.Path)-len(filepath.Ext(file.Path))])
}

func OpenPWKTFile(path string) (*PWKTFile, error) {
	file := PWKTFile{make([][]Token, 0), path}
	bytes, e := ioutil.ReadFile(path)
	if e != nil {
		return &file, errors.New(e.Error())
	}
	data := strings.Replace(string(bytes), "\n\r", "\n", -1)
	for idx, ln := range strings.Split(data, "\n") {
		sl := strings.Trim(ln, "\n\r")
		if len(sl) < 1 {
			continue
		}
		p, e := ParseLine(sl)
		if e == nil {
			file.Lines = append(file.Lines, p)
		} else {
			return &file, errors.New(e.Error() + " at line " + strconv.Itoa(idx))
		}
	}
	return &file, nil
}
