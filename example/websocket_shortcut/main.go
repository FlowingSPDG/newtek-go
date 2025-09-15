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

	// WebSocket経由でのショートカット送信の例
	fmt.Println("\n=== WebSocket経由でのショートカット送信の例 ===")

	// 基本的なショートカット（パラメータなし）
	fmt.Println("1. TakeショートカットをWebSocket経由で送信...")
	if err := c.ShortcutWS(ctx, "take", nil); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// パラメータ付きショートカット
	fmt.Println("2. プレビューソース設定をWebSocket経由で送信...")
	if err := c.ShortcutWS(ctx, "set_preview", map[string]string{"source": "input1"}); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// 複数パラメータ付きショートカット
	fmt.Println("3. M/Eソース設定をWebSocket経由で送信...")
	if err := c.ShortcutWS(ctx, "set_me", map[string]string{"me": "me1", "source": "input2"}); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// DSK制御
	fmt.Println("4. DSKオンエア設定をWebSocket経由で送信...")
	if err := c.ShortcutWS(ctx, "set_dsk_onair", map[string]string{"dsk": "dsk1", "onair": "true"}); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// メディア制御
	fmt.Println("5. メディア再生をWebSocket経由で送信...")
	if err := c.ShortcutWS(ctx, "media_play", map[string]string{"media": "ddr1"}); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 3)

	// メディア停止
	fmt.Println("6. メディア停止をWebSocket経由で送信...")
	if err := c.ShortcutWS(ctx, "media_stop", map[string]string{"media": "ddr1"}); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// 録画制御
	fmt.Println("7. 録画開始をWebSocket経由で送信...")
	if err := c.ShortcutWS(ctx, "start-record", nil); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 3)

	fmt.Println("8. 録画停止をWebSocket経由で送信...")
	if err := c.ShortcutWS(ctx, "stop-record", nil); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// ストリーミング制御
	fmt.Println("9. ストリーミング開始をWebSocket経由で送信...")
	if err := c.ShortcutWS(ctx, "start-stream", nil); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 3)

	fmt.Println("10. ストリーミング停止をWebSocket経由で送信...")
	if err := c.ShortcutWS(ctx, "stop-stream", nil); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// DataLink設定
	fmt.Println("11. DataLink値をWebSocket経由で設定...")
	if err := c.ShortcutWS(ctx, "set_datalink", map[string]string{"key": "title", "value": "WebSocket Test"}); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// マクロ実行
	fmt.Println("12. システムマクロをWebSocket経由で実行...")
	if err := c.ShortcutWS(ctx, "system_macro", map[string]string{"macro": "startup"}); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// 出力制御
	fmt.Println("13. ミックス出力設定をWebSocket経由で送信...")
	if err := c.ShortcutWS(ctx, "set_mix_output", map[string]string{"output": "output1"}); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	fmt.Println("\nWebSocket経由でのショートカット送信の例が完了しました。")
}
