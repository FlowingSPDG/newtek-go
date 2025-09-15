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

	// マクロ実行の例
	fmt.Println("\n=== マクロ実行の例 ===")

	// システムマクロの実行
	fmt.Println("システムマクロ 'startup' を実行...")
	if err := c.RunSystemMacro(ctx, "startup"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// カスタムマクロの実行
	fmt.Println("カスタムマクロ 'custom1' を実行...")
	if err := c.RunCustomMacro(ctx, "custom1"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// DataLink機能の例
	fmt.Println("\n=== DataLink機能の例 ===")

	// DataLink値の設定
	fmt.Println("DataLink値 'title' を 'Hello World' に設定...")
	if err := c.SetDataLinkValue(ctx, "title", "Hello World"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 1)

	// DataLink値の取得
	fmt.Println("DataLink値 'title' を取得...")
	value, err := c.GetDataLinkValue(ctx, "title")
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("取得した値: %s\n", value)
	}

	// 複数のDataLink値を設定
	fmt.Println("\n複数のDataLink値を設定...")
	dataLinks := map[string]string{
		"title":    "NewTek TriCaster",
		"subtitle": "Live Production",
		"time":     "2024-01-01 12:00:00",
		"location": "Studio A",
	}

	for key, value := range dataLinks {
		fmt.Printf("DataLink値 '%s' を '%s' に設定...\n", key, value)
		if err := c.SetDataLinkValue(ctx, key, value); err != nil {
			fmt.Printf("エラー: %v\n", err)
		}
		time.Sleep(time.Millisecond * 500)
	}

	// 設定した値を確認
	fmt.Println("\n設定したDataLink値を確認...")
	for key := range dataLinks {
		value, err := c.GetDataLinkValue(ctx, key)
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
		} else {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	// 出力制御の例
	fmt.Println("\n=== 出力制御の例 ===")
	fmt.Println("ミックス出力を 'output1' に設定...")
	if err := c.SetMixOutput(ctx, "output1"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	fmt.Println("\nマクロ・DataLink・出力制御の例が完了しました。")
}
