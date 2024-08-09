package main

import (
	"fmt"
	api "github.com/quackduck/devzat/devzatapi"
	"sync"
	"time"
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

var sessionsLock sync.Mutex
var instancesSessions []InstanceSession

func makeSessionInstances(insts [2]Instance) {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()
	instancesSessions = make([]InstanceSession, len(insts))
	for i := range insts {
		instancesSessions[i].instance = insts[i]
		session, err := api.NewSession(insts[i].Host, insts[i].Token)
		if err != nil {
			fmt.Println(err)
			instancesSessions[i].connected = false
		} else {
			instancesSessions[i].connected = true
			instancesSessions[i].session = session
		}
	}
}

func courier(msg api.Message, fromInstance int) {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()
	from := instancesSessions[fromInstance].instance.Prefix + " " + msg.From
	for i := range instancesSessions {
		if i == fromInstance || !instancesSessions[i].connected {
			continue
		}
		err := instancesSessions[i].session.SendMessage(api.Message{Room: msg.Room, From: from, Data: msg.Data})
		if err != nil { // TODO: refacto that
			fmt.Println(err)
			instancesSessions[i].connected = false
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
	msg := api.Message{Room: "#main", From: "courier", Data: "Coucou"}
	courier(msg, 1)
	time.Sleep(20 * time.Second)
	courier(msg, 0)
}
