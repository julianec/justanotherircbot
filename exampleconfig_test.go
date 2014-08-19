package main

import (
	"code.google.com/p/goprotobuf/proto"
	"io/ioutil"
	"testing"
)

func TestExampleConfig(t *testing.T) {
	var configdata []byte
	var config IRCBotConfig
	var err error

	configdata, err = ioutil.ReadFile("exampleconfig")
	if err != nil {
		t.Error("Error reading exampleconfig: ", err)
	}

	err = proto.UnmarshalText(string(configdata), &config)
	if err != nil {
		t.Error("Error parsing config file: ", err)
	}
}
