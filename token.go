package main

import (
	"io/ioutil"
	"./errorcheck"
	"gopkg.in/yaml.v2"
)

// T _
type T struct {
	Token string
}

func getToken() string {
	var data, err = ioutil.ReadFile("./token.yaml")
	errorcheck.Check(err)
	t := T{}
	err = yaml.Unmarshal(data, &t)
	errorcheck.Check(err)
	return t.Token
}
