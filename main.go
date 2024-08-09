package main

import (
	//"fmt"
	api "github.com/quackduck/devzat/devzatapi"
)

type Instance struct {
	Host   string
	Token  string
	Prefix string
}

var Instances = [2]Instance{
	Instance{Host: "localhost:5557", Token: "dvz@OB8RwxxDaJzJg2hZclgWuEQD2XkqW1L5zFMpUw7k2gs=", Prefix: "1"},
	Instance{Host: "localhost:5558", Token: "dvz@fX+Rx4eNVuTzfxwKPaQjBZoUksrlDNwMFvQY8A5NhXM=", Prefix: "2"},
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	for i := range Instances {
		session, err := api.NewSession(Instances[i].Host, Instances[i].Token)
		errPanic(err)
		err = session.SendMessage(api.Message{Room: "#main", From: Instances[i].Prefix, Data: "Coucou"})
		errPanic(err)
	}
}
