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

	// 製品情報を取得
	prod, err := c.GetProduct(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Product Model:", prod.ProductModel)
	fmt.Println("Product Name:", prod.ProductName)
	fmt.Println("Session Name:", prod.SessionName)

	// 録画・ストリーミング制御の例
	fmt.Println("\n=== 録画・ストリーミング制御の例 ===")

	// 録画開始
	fmt.Println("録画を開始...")
	if err := c.RecordStart(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 5)

	// ストリーミング開始
	fmt.Println("ストリーミングを開始...")
	if err := c.StreamStart(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 10)

	// ストリーミング停止
	fmt.Println("ストリーミングを停止...")
	if err := c.StreamStop(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 5)

	// 録画停止
	fmt.Println("録画を停止...")
	if err := c.RecordStop(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// トグル機能の例
	fmt.Println("\n=== トグル機能の例 ===")
	fmt.Println("録画をトグル...")
	if err := c.RecordToggle(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 3)

	fmt.Println("ストリーミングをトグル...")
	if err := c.StreamToggle(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 5)

	fmt.Println("録画を再度トグル...")
	if err := c.RecordToggle(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 3)

	fmt.Println("ストリーミングを再度トグル...")
	if err := c.StreamToggle(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	fmt.Println("\n録画・ストリーミング制御の例が完了しました。")
}
