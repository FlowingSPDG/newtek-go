package main

import (
	"encoding/xml"
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

	b, err := xml.Marshal(req)
	if err != nil {
		fmt.Println("Failed to send command:", err)
		return
	}

	if err := c.SendBytes(b); err != nil {
		fmt.Println("Failed to send command:", err)
		return
	}
	return
}
