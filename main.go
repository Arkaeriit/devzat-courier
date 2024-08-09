package main

import (
	"fmt"
	api "github.com/quackduck/devzat/devzatapi"
)

type Instance struct {
	Host   string
	Token  string
	Prefix string
}

type InstanceSession struct {
	instance  Instance
	session   *api.Session
	connected bool
}

var Instances = [2]Instance{
	Instance{Host: "localhost:5557", Token: "dvz@OB8RwxxDaJzJg2hZclgWuEQD2XkqW1L5zFMpUw7k2gs=", Prefix: "1"},
	Instance{Host: "localhost:5558", Token: "dvz@fX+Rx4eNVuTzfxwKPaQjBZoUksrlDNwMFvQY8A5NhXM=", Prefix: "2"},
}

var InstancesSessions []InstanceSession

func makeSessionInstances(insts [2]Instance) {
	InstancesSessions = make([]InstanceSession, len(insts))
	for i := range insts {
		InstancesSessions[i].instance = insts[i]
		session, err := api.NewSession(insts[i].Host, insts[i].Token)
		if err != nil {
			fmt.Println(err)
			InstancesSessions[i].connected = false
		} else {
			InstancesSessions[i].connected = true
			InstancesSessions[i].session = session
		}
	}
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	makeSessionInstances(Instances)
	for i := range InstancesSessions {
		err := InstancesSessions[i].session.SendMessage(api.Message{Room: "#main", From: InstancesSessions[i].instance.Prefix, Data: "Coucou"})
		errPanic(err)
	}
}
