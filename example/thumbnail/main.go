package main

import (
	"fmt"

	"github.com/FlowingSPDG/newtek-go"
)

const (
	host = "192.168.100.136"
)

func main() {
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

	i, err := c.VideoPreview("output1", 1920, 1080)
	if err != nil {
		panic(err)
	}

	fmt.Println("Size:", i.Bounds().Size())
}
