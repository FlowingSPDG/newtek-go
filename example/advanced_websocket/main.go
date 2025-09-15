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

	// 高度なWebSocket機能の例
	fmt.Println("\n=== 高度なWebSocket機能の例 ===")

	// 1. ショートカット状態通知
	fmt.Println("1. ショートカット状態通知を開始...")
	if err := c.ShortcutStateNotifications(ctx, func(msg string) {
		fmt.Printf("ショートカット状態変更: %s\n", msg)
	}); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// 2. オーディオ通知（output1）
	fmt.Println("2. オーディオ通知を開始...")
	if err := c.AudioNotifications(ctx, "output1", func(data []byte) {
		fmt.Printf("オーディオデータ受信: %d bytes\n", len(data))
	}); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// 3. ビデオ通知（output1, 640x480, 品質80）
	fmt.Println("3. ビデオ通知を開始...")
	if err := c.VideoNotifications(ctx, "output1", 640, 480, 80, func(data []byte) {
		fmt.Printf("ビデオフレーム受信: %d bytes\n", len(data))
	}); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// 4. WebSocket経由でショートカットを送信
	fmt.Println("4. WebSocket経由でショートカットを送信...")
	if err := c.ShortcutWS(ctx, "take", nil); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 2)

	// パラメータ付きショートカット
	if err := c.ShortcutWS(ctx, "set_preview", map[string]string{"source": "input1"}); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// 5. 高度なマクロ制御
	fmt.Println("\n=== 高度なマクロ制御の例 ===")

	// マクロの有効性確認
	fmt.Println("マクロ 'startup' の有効性を確認...")
	enabled, err := c.MacroIsEnabledByName(ctx, "startup")
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("マクロ 'startup' は有効: %t\n", enabled)
	}

	// マクロIDの取得
	fmt.Println("マクロ 'startup' のIDを取得...")
	macroID, err := c.GetMacroIDByName(ctx, "startup")
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("マクロID: %s\n", macroID)
	}

	// マクロの記録開始
	fmt.Println("マクロ 'test_macro' の記録を開始...")
	if err := c.RecordMacroByName(ctx, "test_macro"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	time.Sleep(time.Second * 3)

	// マクロの記録停止
	fmt.Println("マクロの記録を停止...")
	if err := c.StopRecordMacro(ctx); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// 記録したマクロを実行
	fmt.Println("記録したマクロを実行...")
	if err := c.PlayMacroByName(ctx, "test_macro"); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}

	// 6. 辞書クエリ機能
	fmt.Println("\n=== 辞書クエリ機能の例 ===")

	// スイッチャー情報の取得
	fmt.Println("スイッチャー情報を取得...")
	switcherInfo, err := c.GetSwitcherInfo(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("スイッチャー情報: %s\n", string(switcherInfo)[:min(200, len(switcherInfo))])
	}

	// バッファ情報の取得
	fmt.Println("バッファ情報を取得...")
	bufferInfo, err := c.GetBufferInfo(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("バッファ情報: %s\n", string(bufferInfo)[:min(200, len(bufferInfo))])
	}

	// マクロリストの取得
	fmt.Println("マクロリストを取得...")
	macrosList, err := c.GetMacrosList(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("マクロリスト: %s\n", string(macrosList)[:min(200, len(macrosList))])
	}

	// NDIソース情報の取得
	fmt.Println("NDIソース情報を取得...")
	ndiSources, err := c.GetNDISources(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("NDIソース情報: %s\n", string(ndiSources)[:min(200, len(ndiSources))])
	}

	// DDRプレイリストの取得
	fmt.Println("DDRプレイリストを取得...")
	ddrPlaylist, err := c.GetDDRPlaylist(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("DDRプレイリスト: %s\n", string(ddrPlaylist)[:min(200, len(ddrPlaylist))])
	}

	// DDRタイムコードの取得
	fmt.Println("DDRタイムコードを取得...")
	ddrTimecode, err := c.GetDDRTimecode(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("DDRタイムコード: %s\n", string(ddrTimecode)[:min(200, len(ddrTimecode))])
	}

	// オーディオミキサー情報の取得
	fmt.Println("オーディオミキサー情報を取得...")
	audioMixer, err := c.GetAudioMixerInfo(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("オーディオミキサー情報: %s\n", string(audioMixer)[:min(200, len(audioMixer))])
	}

	// ファイルブラウザ情報の取得
	fmt.Println("ファイルブラウザ情報を取得...")
	fileBrowser, err := c.GetFileBrowserInfo(ctx)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Printf("ファイルブラウザ情報: %s\n", string(fileBrowser)[:min(200, len(fileBrowser))])
	}

	fmt.Println("\n高度なWebSocket機能の例が完了しました。")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
