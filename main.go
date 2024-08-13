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

/* ---------------------------------- Main ---------------------------------- */

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

/* ----------------------------- Initialization ----------------------------- */

// Takes a list of template instances and use it to open all the sessions.
func makeSessionInstances(insts []Instance) {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()
	instancesSessions = make([]InstanceSession, len(insts))
	messagesChan = make(chan MessageFrom, len(insts)*2)
	for i := range insts {
		instancesSessions[i].openSession(insts[i])
		instancesSessions[i].registerListener()
		instancesSessions[i].registerCmd()
	}
}

// Start the initialization of the InstanceSession by copying the template and
// opening the session.
func (instance *InstanceSession) openSession(template Instance) {
	instance.instance = template
	session, err := api.NewSession(template.Host, template.Token)
	instance.connected = true
	instance.session = session
	manageInstanceError(instance, err)
}

// Register a wildcard listener for every session.
func (instance *InstanceSession) registerListener() {
	if !instance.connected {
		return
	}
	msgChan, _, err := instance.session.RegisterListener(false, false, "")
	instance.msgChan = msgChan
	manageInstanceError(instance, err)
}

// Register the courier command to every session.
func (instance *InstanceSession) registerCmd() {
	if !instance.connected {
		return
	}
	err := instance.session.RegisterCmd("courier", "command", "Run 'courier help' for more information.",
		func(cmdCall api.CmdCall, err error) {
			sessionsLock.Lock()
			defer sessionsLock.Unlock()
			courierCmd(instance.session, cmdCall, instancesSessions)
		})
	manageInstanceError(instance, err)
}

/* -------------------------- Message transmission -------------------------- */

// Periodically reads messages and dispatch them. Sleep between runs to reduce
// CPU usage.
func courierLoop() {
	for {
		readMsgChans()
		dispatchMessages()
		time.Sleep(time.Millisecond * 250)
	}
}

// Reads messages from every instances and push them to the queue.
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

// Pop messages from the queue and distribute them to the instances.
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

// Takes a message and distribute it to every instance but the instance it comes
// from. Add needed prefix and colors to the message.
func courier(msg MessageFrom) {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()
	prefix := colorPrefix(instancesSessions[msg.fromInstance].instance)
	user := colorName(instancesSessions[msg.fromInstance].instance, msg.msg.From)
	from := prefix + " " + user
	for i := range instancesSessions {
		if i == msg.fromInstance || !instancesSessions[i].connected {
			continue
		}
		err := instancesSessions[i].session.SendMessage(api.Message{Room: msg.msg.Room, From: from, Data: msg.msg.Data})
		manageInstanceError(&instancesSessions[i], err)
	}
}

/* --------------------------------- Helpers -------------------------------- */

// Returns the prefix of an instance with the proper color.
func colorPrefix(i Instance) string {
	return colorString(i.Prefix, i.PrefixColor)
}

// Returns `name` colored as a name according to the instance's configuration.
func colorName(i Instance, name string) string {
	return colorString(name, i.NameColor)
}

// Takes an error and an instance, if the error is nil or acceptable, do
// nothing. Otherwise, prints it and disconnect the instance.
func manageInstanceError(instance *InstanceSession, err error) {
	if err == nil {
		return
	}
	if err.Error() == "rpc error: code = InvalidArgument desc = Room does not exist" {
		// This is what happen if you talk in a room which doesn't exists in
		// other instances, totally normal.
		return
	}
	fmt.Println(err.Error())
	fmt.Println(err)
	instance.connected = false
}

// If the error is not nil, panic with it.
func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}
