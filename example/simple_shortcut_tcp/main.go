package main

import (
	"fmt"

	"github.com/FlowingSPDG/newtek-go"
)

const (
	host = "192.168.100.93"
)

func main() {
	c, err := newtek.NewClientV1TCP(host)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	shortcuts := make([]newtek.Shortcut, 0)
	shortcuts = append(shortcuts, newtek.Shortcut{
		Name:  "mode",
		Value: "2",
	})

	req := newtek.Shortcuts{
		Shortcut: shortcuts,
	}
	if err := c.Shortcut(req); err != nil {
		fmt.Println("Failed to send shortcut:", err)
		return // DO NOT panic because TCP connection needs to be closed
	}

	return
}
