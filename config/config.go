package config

import (
	"io/ioutil"

	"../errorcheck"
	"gopkg.in/yaml.v2"
)

// Config _
type Config struct {
	Token  string
	Prefix string
	Status string
}

func GetConfig() Config {
	var data, err = ioutil.ReadFile("./config.yaml")
	errorcheck.Check(err)
	c := Config{}
	err = yaml.Unmarshal(data, &c)
	errorcheck.Check(err)
	return c
}
