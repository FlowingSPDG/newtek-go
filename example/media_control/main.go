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

	// メディア制御の例
	fmt.Println("\n=== メディア制御の例 ===")

	// DDR1のメディアを再生
	fmt.Println("DDR1のメディアを再生...")
	if err := c.MediaPlay(ctx, "ddr1"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 3)

	// 一時停止
	fmt.Println("メディアを一時停止...")
	if err := c.MediaPause(ctx, "ddr1"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// 再開
	fmt.Println("メディアを再開...")
	if err := c.MediaPlay(ctx, "ddr1"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 3)

	// 巻き戻し
	fmt.Println("メディアを巻き戻し...")
	if err := c.MediaRewind(ctx, "ddr1"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// 早送り
	fmt.Println("メディアを早送り...")
	if err := c.MediaFastForward(ctx, "ddr1"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// 停止
	fmt.Println("メディアを停止...")
	if err := c.MediaStop(ctx, "ddr1"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// GFXメディアの制御例
	fmt.Println("\n=== GFXメディアの制御例 ===")
	fmt.Println("GFX1のメディアを再生...")
	if err := c.MediaPlay(ctx, "gfx1"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 3)

	fmt.Println("GFX1のメディアを停止...")
	if err := c.MediaStop(ctx, "gfx1"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	fmt.Println("\nメディア制御の例が完了しました。")
}
