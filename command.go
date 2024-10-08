package main

import (
	"fmt"
	"regexp"

	api "github.com/quackduck/devzat/devzatapi"
)

/* ---------------------------------- Main ---------------------------------- */

func courierCmd(session *api.Session, cmdCall api.CmdCall, instances []InstanceSession) {
	var msg string
	if cmdCall.Args == "help" {
		msg = help
	} else if cmdCall.Args == "status" {
		msg = formatStatus(instances)
	} else {
		return
	}
	session.SendMessage(api.Message{Room: cmdCall.Room, From: "", DMTo: "", Data: msg})
}

/* ---------------------------------- Help ---------------------------------- */

const help = `Devzat Courier, a plugin to let multiple Devzat instance communicate.

Each message will be sent to every connected instances with an identifying prefix.

Source code available on [GitHub](https://github.com/Arkaeriit/devzat-courier)
` +
	"\nTo list the connected instances and their state, do `courier status`.\n"

/* --------------------------------- Status --------------------------------- */

// Return the host without the port of an instance.
func getHost(i Instance) string {
	r, _ := regexp.Compile(":[0-9]+$")
	ret := r.ReplaceAllString(i.Host, "")
	return ret
}

// Looks like someone is making code very coupled with the rest of the system...
func longuestHostName(instances []InstanceSession) int {
	ret := 0
	for _, i := range instances {
		hostLen := len(colorName(i.instance, getHost(i.instance)))
		if hostLen > ret {
			ret = hostLen
		}
	}
	return ret
}
func longuestPrefix(instances []InstanceSession) int {
	ret := 0
	for _, i := range instances {
		prefixLen := len(colorPrefix(i.instance))
		if prefixLen > ret {
			ret = prefixLen
		}
	}
	return ret
}

// Makes a status line from an instance session, with pretty host and prefix and
// connected state.
func formatInstanceStatus(i InstanceSession, prefixLen, hostLen int) string {
	prefix := colorPrefix(i.instance)
	host := colorName(i.instance, getHost(i.instance))
	status := colorString("online", "green")
	if !i.connected {
		status = colorString("offline", "red")
	}
	return fmt.Sprintf("%- *s %- *s %s\n", prefixLen, prefix, hostLen, host, status)
}

// Makes a status line for every instances.
func formatStatus(instances []InstanceSession) string {
	ret := ""
	prefixLen := longuestPrefix(instances)
	hostLen := longuestHostName(instances)
	for _, i := range instances {
		ret += formatInstanceStatus(i, prefixLen, hostLen)
	}
	return ret
}
