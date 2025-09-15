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

	// フィードバック・状態監視の例
	fmt.Println("\n=== フィードバック・状態監視の例 ===")

	// タリー情報の取得
	fmt.Println("タリー情報を取得...")
	tally, err := c.GetTally(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("タリー情報: %+v\n", tally)
		for _, column := range tally.Column {
			fmt.Printf("  ソース: %s, プログラム: %s, プレビュー: %s\n",
				column.Name, column.OnPgm, column.OnPrev)
		}
	}

	// 録画状態の確認
	fmt.Println("\n録画状態を確認...")
	recording, err := c.GetRecordingStatus(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("録画状態: %t\n", recording)
	}

	// ストリーミング状態の確認
	fmt.Println("ストリーミング状態を確認...")
	streaming, err := c.GetStreamingStatus(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("ストリーミング状態: %t\n", streaming)
	}

	// ソース状態の確認
	fmt.Println("\nソース状態を確認...")
	program, preview, err := c.GetSourceStatus(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("プログラムソース: %s\n", program)
		fmt.Printf("プレビューソース: %s\n", preview)
	}

	// メディア状態の確認
	fmt.Println("\nメディア状態を確認...")
	mediaList := []string{"ddr1", "ddr2", "gfx1", "gfx2", "stills1", "titles1"}
	for _, media := range mediaList {
		status, err := c.GetMediaStatus(ctx, media)
		if err != nil {
			fmt.Printf("エラー (%s): %v\n", media, err)
		} else {
			fmt.Printf("%s の状態: %s\n", media, status)
		}
	}

	// リアルタイム状態監視の例
	fmt.Println("\n=== リアルタイム状態監視の例 ===")
	fmt.Println("10秒間、状態を監視します...")

	for i := 0; i < 10; i++ {
		fmt.Printf("\n--- 監視 %d/10 ---\n", i+1)

		// 録画状態
		recording, err := c.GetRecordingStatus(ctx)
		if err != nil {
			fmt.Printf("録画状態取得エラー: %v\n", err)
		} else {
			fmt.Printf("録画状態: %t\n", recording)
		}

		// ストリーミング状態
		streaming, err := c.GetStreamingStatus(ctx)
		if err != nil {
			fmt.Printf("ストリーミング状態取得エラー: %v\n", err)
		} else {
			fmt.Printf("ストリーミング状態: %t\n", streaming)
		}

		// ソース状態
		program, preview, err := c.GetSourceStatus(ctx)
		if err != nil {
			fmt.Printf("ソース状態取得エラー: %v\n", err)
		} else {
			fmt.Printf("プログラム: %s, プレビュー: %s\n", program, preview)
		}

		time.Sleep(time.Second * 1)
	}

	// ショートカット状態の確認
	fmt.Println("\n=== ショートカット状態の確認 ===")
	shortcutStates, err := c.ShortcutStates(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("ショートカット状態数: %d\n", len(shortcutStates.ShortcutState))
		for i, state := range shortcutStates.ShortcutState {
			if i < 10 { // 最初の10個のみ表示
				fmt.Printf("  %s: %s (タイプ: %s)\n",
					state.Name, state.Value, state.Type)
			}
		}
		if len(shortcutStates.ShortcutState) > 10 {
			fmt.Printf("  ... 他 %d 個のショートカット状態\n",
				len(shortcutStates.ShortcutState)-10)
		}
	}

	fmt.Println("\nフィードバック・状態監視の例が完了しました。")
}
