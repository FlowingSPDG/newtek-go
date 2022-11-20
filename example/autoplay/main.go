package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/FlowingSPDG/newtek-go"
)

var (
	tpHost = "" // 3Play host
	tcHost = "" // TriCaster host
)

func main() {
	// Parse flags
	flag.StringVar(&tpHost, "3play", "10.36.11.8", "3Play IP Address")
	flag.StringVar(&tcHost, "tricaster", "10.36.11.6", "TriCaster IP Address")
	flag.Parse()

	// Solve 3Play
	tp, err := newtek.NewClientV1(tpHost, "admin", "admin")
	if err != nil {
		panic(err)
	}
	prodTP, err := tp.GetProduct()
	if err != nil {
		panic(err)
	}
	fmt.Println("[3Play] Product Model:", prodTP.ProductModel)
	fmt.Println("[3Play] Product Name:", prodTP.ProductName)
	fmt.Println("[3Play] Session Name:", prodTP.SessionName)

	// Solve TriCsater
	tc, err := newtek.NewClientV1(tcHost, "admin", "admin")
	if err != nil {
		panic(err)
	}
	prodTC, err := tc.GetProduct()
	if err != nil {
		panic(err)
	}
	fmt.Println("[TriCaster] Product Model:", prodTC.ProductModel)
	fmt.Println("[TriCaster] Product Name:", prodTC.ProductName)
	fmt.Println("[TriCaster] Session Name:", prodTC.SessionName)

	// [3Play] クリップを取得
	m, err := tp.MetaData()
	if err != nil {
		panic(err)
	}

	// [3Play] Enter cliplist mode
	if err := tp.ShortcutHTTP("mode", map[string]string{"value": "0"}); err != nil {
		panic(err)
	}

	// まれに最後のクリップが追加されないことがある
	// [3Play] Only use first clip list for this time
	for i := 0; i < len(m.Clips[0].Mark); i++ {
		mark := m.Clips[0].Mark[i]

		// [3Play] Select EventID and angle(0)
		if err := tp.ShortcutHTTP("out1_clip_select", map[string]string{"value": fmt.Sprintf("%s-0", mark.ID)}); err != nil {
			panic(err)
		}

		// [3Play] Add to playlist
		if err := tp.ShortcutHTTP("add-to-list", nil); err != nil {
			panic(err)
		}

		// Delay
		time.Sleep(time.Millisecond * 100)
	}

	time.Sleep(time.Millisecond * 250)

	// [3Play] Switch to Playlist mode
	if err := tp.ShortcutHTTP("mode", map[string]string{"value": "1"}); err != nil {
		panic(err)
	}

	// [TriCaster] Ready preview(b_row) for M/E1(v1). ちょっと遅延がある?
	if err := tc.ShortcutHTTP("main_b_row_named_input", map[string]string{"input": "v1"}); err != nil {
		panic(err)
	}

	// Delay
	time.Sleep(time.Millisecond * 250)

	// [3Play] Play
	if err := tp.ShortcutHTTP("play", nil); err != nil {
		panic(err)
	}

	// [TriCaster] Stinger into 3Play with Auto
	if err := tc.ShortcutHTTP("main_auto", nil); err != nil {
		panic(err)
	}

	// 自動削除も欲しい
}
