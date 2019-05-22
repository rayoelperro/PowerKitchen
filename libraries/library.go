package libraries

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Method struct {
	Name string
	Code string
}

type Library struct {
	Name         string
	Load         string
	Var          string
	Standard     bool
	Methods      []Method
	Properties   []string
	Requirements []string
	Info         []string
}

type LibCollection struct {
	Libraries []Library
}

func Libs(instructions string) (*LibCollection, error) {
	result := make([]Library, 0)
	content, e := ioutil.ReadFile(instructions)
	if e != nil {
		return &LibCollection{result}, e
	}
	json.Unmarshal(content, &result)
	return &LibCollection{result}, nil
}

func (lc *LibCollection) Take(name string) (*Library, error) {
	for _, l := range lc.Libraries {
		if l.Name == name {
			return &l, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("The library %s does not exists", name))
}

func (lb *Library) MethodsName() []string {
	methods := make([]string, 0)
	for _, v := range lb.Methods {
		methods = append(methods, v.Name)
	}
	return methods
}
