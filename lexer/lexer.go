package lexer

import (
	"errors"
	"strconv"
	"strings"
)

type TokenClass string

const (
	Head TokenClass = "Head"
	Arg  TokenClass = "Arg"
)

type Token struct {
	Class TokenClass
	Value string
}

func GenToken(str string, idx int) Token {
	if idx == 0 {
		return Token{Head, str}
	} else {
		return Token{Arg, str}
	}
}

func ParseLine(line string) ([]Token, error) {
	tks := make([]Token, 0)
	actual := ""
	inmod := false
	instr := false
	is := 0
	for i, c := range line {
		if (i == len(line)-1) || (c == '"' && !inmod) || (c == ' ' && !instr) {
			if c == '\\' && !inmod {
				return tks, errors.New("You forgot to close a character shortcut")
			}
			wasstr := false
			if c == '"' && !inmod {
				wasstr = instr
				instr = !instr
			}
			if i == len(line)-1 {
				if instr {
					return tks, errors.New("You forgot to close an string")
				} else {
					if !wasstr {
						actual += string(c)
					}
				}
			}
			if len(actual) > 0 || wasstr {
				tks = append(tks, GenToken(actual, is))
			}
			actual = ""
		} else {
			if actual == "" {
				is = i
			}
			if c == '\\' && !inmod {
				inmod = true
			} else {
				actual += string(c)
				inmod = false
			}
		}
	}
	return tks, nil
}

func MakeString(tk []Token) string {
	result := ""
	for i, e := range tk {
		wll := strings.Replace(strings.Replace(strings.Replace(
			e.Value, "\\", "\\\\", -1), "\"", "\\\"", -1), "\n", "\\n", -1)
		if contains(' ', e.Value) {
			result += "\"" + wll + "\""
		} else {
			result += wll
		}
		if i < len(tk)-1 {
			result += " "
		}
	}
	return result
}

func contains(e rune, g string) bool {
	for _, p := range g {
		if p == e {
			return true
		}
	}
	return false
}

func AllowLength(tk []Token, ln int, fn func() error) error {
	if ln == len(tk) {
		return fn()
	} else if ln > len(tk) {
		return errors.New("You send " + strconv.Itoa(ln-len(tk)) + " less arguments")
	}
	return errors.New("You send " + strconv.Itoa(len(tk)-ln) + " more arguments")
}

func AllowAtLeastLength(tk []Token, ln int, fn func() error) error {
	if len(tk) >= ln {
		return fn()
	}
	return errors.New("You send " + strconv.Itoa(len(tk)-ln) + " more arguments")
}

func ToStringArray(tk []Token) []string {
	sf := make([]string, 0)
	for _, t := range tk {
		sf = append(sf, t.Value)
	}
	return sf
}
