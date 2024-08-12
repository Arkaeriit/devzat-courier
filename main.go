package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	api "github.com/quackduck/devzat/devzatapi"
)

type Instance struct {
	Host        string
	Token       string
	Prefix      string
	NameColor   string
	PrefixColor string
}

type InstanceSession struct {
	instance  Instance
	session   *api.Session
	msgChan   chan api.Message
	connected bool
}

type MessageFrom struct {
	msg          api.Message
	fromInstance int
}

var sessionsLock sync.Mutex
var instancesSessions []InstanceSession
var messagesChan chan MessageFrom

func makeSessionInstances(insts []Instance) {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()
	instancesSessions = make([]InstanceSession, len(insts))
	messagesChan = make(chan MessageFrom, len(insts)*2)
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
		msgChan, _, err := session.RegisterListener(false, false, "")
		instancesSessions[i].msgChan = msgChan
		if err != nil { // TODO: refacto that
			fmt.Println(err)
			instancesSessions[i].connected = false
		}
	}
}

func courier(msg MessageFrom) {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()
	prefix := colorString(instancesSessions[msg.fromInstance].instance.Prefix, instancesSessions[msg.fromInstance].instance.PrefixColor)
	user := colorString(msg.msg.From, instancesSessions[msg.fromInstance].instance.NameColor)
	from := prefix + " " + user
	for i := range instancesSessions {
		if i == msg.fromInstance || !instancesSessions[i].connected {
			continue
		}
		err := instancesSessions[i].session.SendMessage(api.Message{Room: msg.msg.Room, From: from, Data: msg.msg.Data})
		if err != nil { // TODO: refacto that
			fmt.Println(err)
			instancesSessions[i].connected = false
		}
	}
}

func readMsgChans() {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()
	for i := range instancesSessions {
		if !instancesSessions[i].connected {
			continue
		}
		select {
		case err := <-instancesSessions[i].session.ErrorChan:
			fmt.Println(err)
			instancesSessions[i].connected = false
		case msg := <-instancesSessions[i].msgChan:
			msgFrom := MessageFrom{msg: msg, fromInstance: i}
			messagesChan <- msgFrom
		default:
			continue
		}
	}
}

func dispatchMessages() {
	for {
		select {
		case msg := <-messagesChan:
			courier(msg)
		default:
			return
		}
	}
}

func courierLoop() {
	for {
		readMsgChans()
		dispatchMessages()
		time.Sleep(time.Millisecond * 250)
	}
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <config file>\n", os.Args[0])
		os.Exit(1)
	}
	configFile := os.Args[1]
	f, err := os.Open(configFile)
	errPanic(err)
	var instances []Instance
	err = json.NewDecoder(f).Decode(&instances)
	errPanic(err)
	makeSessionInstances(instances)
	fmt.Println("Starting loop")
	courierLoop()
}
