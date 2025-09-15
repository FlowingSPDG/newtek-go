package main

import (
	"context"
	"fmt"
	"time"

	"github.com/FlowingSPDG/newtek-go"
)

const (
	host = "192.168.100.136"
)

func main() {
	ctx := context.Background()
	c, err := newtek.NewClientV1(host, "admin", "admin")
	if err != nil {
		panic(err)
	}

	prod, err := c.GetProduct(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Product Model:", prod.ProductModel)
	fmt.Println("Product Name:", prod.ProductName)
	fmt.Println("Session Name:", prod.SessionName)

	if err := c.ShortcutHTTP(ctx, "record", nil); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * time.Duration(3))

	if err := c.ShortcutHTTP(ctx, "stop-record", nil); err != nil {
		panic(err)
	}
}
