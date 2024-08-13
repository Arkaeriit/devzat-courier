package main

import "fmt"

var colors = map[string]string{
	"black":  "\033[0;30m",
	"red":    "\033[0;31m",
	"green":  "\033[0;32m",
	"yellow": "\033[0;33m",
	"blue":   "\033[0;34m",
	"purple": "\033[0;35m",
	"cyan":   "\033[0;36m",
	"white":  "\033[0;37m",
}

const nc = "\033[m"

// Try to color `in` in the given color. If the color is invalid, return it
// unchanged.
func colorString(in string, color string) string {
	color, ok := colors[color]
	if !ok {
		fmt.Printf("Error, invalid color '%v', no color applied.\n", color)
		return in
	}
	ret := fmt.Sprintf("%v%v%v", color, in, nc)
	return ret
}
