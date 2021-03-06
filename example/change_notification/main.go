package main

import (
	"fmt"

	"github.com/FlowingSPDG/newtek-go"
)

const (
	host = "192.168.100.136"
)

func main() {
	forever := make(chan struct{})
	c, err := newtek.NewClientV1(host, "admin", "admin")
	if err != nil {
		panic(err)
	}

	prod, err := c.GetProduct()
	if err != nil {
		panic(err)
	}
	fmt.Println("Product Model:", prod.ProductModel)
	fmt.Println("Product Name:", prod.ProductName)
	fmt.Println("Session Name:", prod.SessionName)

	if err := c.ChangeNotifications(func(msg string) {
		fmt.Println("msg:", msg)
	}); err != nil {
		panic(err)
	}
	<-forever
}
