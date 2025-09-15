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

	// 基本的なスイッチャー制御の例
	fmt.Println("\n=== スイッチャー制御の例 ===")

	// プレビューにソースを設定
	fmt.Println("プレビューにソース 'input1' を設定...")
	if err := c.SetSourceToPreview(ctx, "input1"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// Take（プログラムとプレビューの切り替え）
	fmt.Println("Take実行...")
	if err := c.Take(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// プレビューに別のソースを設定
	fmt.Println("プレビューにソース 'input2' を設定...")
	if err := c.SetSourceToPreview(ctx, "input2"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// 自動トランジション
	fmt.Println("自動トランジション実行...")
	if err := c.AutoTransition(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// M/Eにソースを設定
	fmt.Println("M/E1にソース 'input3' を設定...")
	if err := c.SetSourceToME(ctx, "me1", "input3"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// DSKをオンエアに設定
	fmt.Println("DSK1をオンエアに設定...")
	if err := c.SetDSKOnAir(ctx, "dsk1", true); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// トランジション選択
	fmt.Println("トランジション 'mix' を選択...")
	if err := c.SetTransitionSelection(ctx, "mix"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	fmt.Println("\nスイッチャー制御の例が完了しました。")
}
