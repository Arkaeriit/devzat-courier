package main

import (
	"fmt"
	"regexp"

	api "github.com/quackduck/devzat/devzatapi"
)

func getHost(i Instance) string {
	r, _ := regexp.Compile(":[0-9]+$")
	ret := r.ReplaceAllString(i.Host, "")
	return ret
}

// Looks like someone is making code very coupled with the rest of the system...
func longuestHostName(instances []InstanceSession) int {
	ret := 0
	for _, i := range instances {
		hostLen := len(getHost(i.instance))
		if hostLen > ret {
			ret = hostLen
		}
	}
	return ret
}
func longuestPrefix(instances []InstanceSession) int {
	ret := 0
	for _, i := range instances {
		prefixLen := len(i.instance.Prefix)
		if prefixLen > ret {
			ret = prefixLen
		}
	}
	return ret
}

func formatInstanceStatus(i InstanceSession, prefixLen, hostLen int) string {
	prefix := colorPrefix(i.instance)
	host := colorName(i.instance, i.instance.Host)
	status := colorString("online", "green")
	if !i.connected {
		status = colorString("offline", "red")
	}
	return fmt.Sprintf("%- *s %- *s %s\n", prefixLen, prefix, hostLen, host, status)
}

func formatStatus(instances []InstanceSession) string {
	ret := ""
	prefixLen := longuestPrefix(instances)
	hostLen := longuestHostName(instances)
	for _, i := range instances {
		ret += formatInstanceStatus(i, prefixLen, hostLen)
	}
	return ret
}

func courierCmd(session *api.Session, cmdCall api.CmdCall, instances []InstanceSession) {
	session.SendMessage(api.Message{Room: cmdCall.Room, From: "", DMTo: "", Data: formatStatus(instances)})
}
