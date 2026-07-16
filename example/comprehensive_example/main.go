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
	fmt.Println("=== NewTek TriCaster 接続情報 ===")
	fmt.Printf("製品モデル: %s\n", prod.ProductModel)
	fmt.Printf("製品名: %s\n", prod.ProductName)
	fmt.Printf("セッション名: %s\n", prod.SessionName)
	fmt.Printf("解像度: %sx%s\n", prod.SessionXResolution, prod.SessionYResolution)
	fmt.Printf("フレームレート: %s fps\n", prod.SessionFrameRate)

	// ライブ状態の確認
	fmt.Println("\n=== ライブ状態の確認 ===")
	live := c.GetLive(ctx)
	fmt.Printf("ライブ状態: %t\n", live)

	// 現在の状態を確認
	fmt.Println("\n=== 現在の状態確認 ===")

	// 録画・ストリーミング状態
	recording, _ := c.GetRecordingStatus(ctx)
	streaming, _ := c.GetStreamingStatus(ctx)
	fmt.Printf("録画状態: %t, ストリーミング状態: %t\n", recording, streaming)

	// ソース状態
	program, preview, _ := c.GetSourceStatus(ctx)
	fmt.Printf("プログラムソース: %s, プレビューソース: %s\n", program, preview)

	// 基本的な放送制御のデモ
	fmt.Println("\n=== 基本的な放送制御デモ ===")

	// 1. プレビューにソースを設定
	fmt.Println("1. プレビューにソース 'input1' を設定...")
	if err := c.SetSourceToPreview(ctx, "input1"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}
	time.Sleep(time.Second * 2)

	// 2. Take実行
	fmt.Println("2. Take実行...")
	if err := c.Take(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}
	time.Sleep(time.Second * 2)

	// 3. 別のソースをプレビューに設定
	fmt.Println("3. プレビューにソース 'input2' を設定...")
	if err := c.SetSourceToPreview(ctx, "input2"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}
	time.Sleep(time.Second * 2)

	// 4. 自動トランジション
	fmt.Println("4. 自動トランジション実行...")
	if err := c.AutoTransition(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}
	time.Sleep(time.Second * 2)

	// メディア制御のデモ
	fmt.Println("\n=== メディア制御デモ ===")

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

	// 停止
	fmt.Println("メディアを停止...")
	if err := c.MediaStop(ctx, "ddr1"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// 録画・ストリーミング制御のデモ
	fmt.Println("\n=== 録画・ストリーミング制御デモ ===")

	// 録画開始
	fmt.Println("録画を開始...")
	if err := c.RecordStart(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}
	time.Sleep(time.Second * 3)

	// ストリーミング開始
	fmt.Println("ストリーミングを開始...")
	if err := c.StreamStart(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}
	time.Sleep(time.Second * 5)

	// 状態確認
	recording, _ = c.GetRecordingStatus(ctx)
	streaming, _ = c.GetStreamingStatus(ctx)
	fmt.Printf("現在の状態 - 録画: %t, ストリーミング: %t\n", recording, streaming)

	// ストリーミング停止
	fmt.Println("ストリーミングを停止...")
	if err := c.StreamStop(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}
	time.Sleep(time.Second * 2)

	// 録画停止
	fmt.Println("録画を停止...")
	if err := c.RecordStop(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// DataLink機能のデモ
	fmt.Println("\n=== DataLink機能デモ ===")

	// 複数のDataLink値を設定
	dataLinks := map[string]string{
		"title":    "NewTek TriCaster Demo",
		"subtitle": "Go言語による制御",
		"time":     time.Now().Format("2006-01-02 15:04:05"),
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

	// マクロ実行のデモ
	fmt.Println("\n=== マクロ実行デモ ===")

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

	// 最終状態確認
	fmt.Println("\n=== 最終状態確認 ===")

	// 録画・ストリーミング状態
	recording, _ = c.GetRecordingStatus(ctx)
	streaming, _ = c.GetStreamingStatus(ctx)
	fmt.Printf("録画状態: %t, ストリーミング状態: %t\n", recording, streaming)

	// ソース状態
	program, preview, _ = c.GetSourceStatus(ctx)
	fmt.Printf("プログラムソース: %s, プレビューソース: %s\n", program, preview)

	// タリー情報
	tally, err := c.GetTally(ctx)
	if err == nil {
		fmt.Println("タリー情報:")
		for _, column := range tally.Column {
			fmt.Printf("  %s: プログラム=%s, プレビュー=%s\n",
				column.Name, column.OnPgm, column.OnPrev)
		}
	}

	fmt.Println("\n=== 包括的なデモが完了しました ===")
}
