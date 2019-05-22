package compiler

import (
	"errors"

	"../remotecontrol"

	. "../lexer"
)

func Run(tk []Token, c *Console) error {
	if len(tk) > 0 {
		if tk[0].Class == Head {
			if isInside(tk[0].Value, Commands) {
				return c.Current.Enter(tk)
			} else if isInside(tk[0].Value, remotecontrol.Commands) {
				return c.Network.Command(tk)
			} else if isInside(tk[0].Value, GlobalCommands) {
				return c.Command(tk)
			}
			return errors.New("Not recognized command: " + tk[0].Value)
		}
		return errors.New("The line had no head")
	}
	return nil
}

func isInside(value string, group []string) bool {
	for _, v := range group {
		if v == value {
			return true
		}
	}
	return false
}
