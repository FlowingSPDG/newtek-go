package newtek

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/icholy/digest"
)

type ClientV1 interface {
	GetLive(ctx context.Context) bool
	GetProduct(ctx context.Context) (info *ProductInformation, err error)
	ShortcutHTTP(ctx context.Context, name string, kv map[string]string) error
	ShortcutWS(ctx context.Context, name string, kv map[string]string) error
	Trigger(ctx context.Context, name string) error

	// Basic Switcher Controls
	Take(ctx context.Context) error
	AutoTransition(ctx context.Context) error
	SetSourceToPreview(ctx context.Context, source string) error
	SetSourceToProgram(ctx context.Context, source string) error
	SetSourceToME(ctx context.Context, me, source string) error
	SetDSKOnAir(ctx context.Context, dsk string, onAir bool) error
	SetTransitionSelection(ctx context.Context, transition string) error

	// Media Controls
	MediaPlay(ctx context.Context, media string) error
	MediaStop(ctx context.Context, media string) error
	MediaRewind(ctx context.Context, media string) error
	MediaFastForward(ctx context.Context, media string) error
	MediaPause(ctx context.Context, media string) error

	// Recording and Streaming
	RecordToggle(ctx context.Context) error
	RecordStart(ctx context.Context) error
	RecordStop(ctx context.Context) error
	StreamToggle(ctx context.Context) error
	StreamStart(ctx context.Context) error
	StreamStop(ctx context.Context) error

	// Macros
	RunSystemMacro(ctx context.Context, macro string) error
	RunCustomMacro(ctx context.Context, macro string) error

	// Advanced Macro Controls
	PlayMacroByName(ctx context.Context, name string) error
	PlayMacroByID(ctx context.Context, id string) error
	RecordMacroByName(ctx context.Context, name string) error
	StopRecordMacro(ctx context.Context) error
	StopMacroByName(ctx context.Context, name string) error
	StopMacroByID(ctx context.Context, id string) error
	MacroIsEnabledByName(ctx context.Context, name string) (bool, error)
	MacroIsEnabledByID(ctx context.Context, id string) (bool, error)
	GetMacroIDByName(ctx context.Context, name string) (string, error)

	// Output Control
	SetMixOutput(ctx context.Context, output string) error

	// DataLink
	SetDataLinkValue(ctx context.Context, key, value string) error
	GetDataLinkValue(ctx context.Context, key string) (string, error)

	// Feedback and Status
	GetTally(ctx context.Context) (*Tally, error)
	GetRecordingStatus(ctx context.Context) (bool, error)
	GetStreamingStatus(ctx context.Context) (bool, error)
	GetMediaStatus(ctx context.Context, media string) (string, error)
	GetSourceStatus(ctx context.Context) (program, preview string, err error)

	// Dictionary Queries
	GetDictionary(ctx context.Context, key string) ([]byte, error)
	GetSwitcherInfo(ctx context.Context) ([]byte, error)
	GetBufferInfo(ctx context.Context) ([]byte, error)
	GetSwitcherUIEffects(ctx context.Context) ([]byte, error)
	GetAudioMixerInfo(ctx context.Context) ([]byte, error)
	GetAudioBinsInfo(ctx context.Context) ([]byte, error)
	GetFileBrowserInfo(ctx context.Context) ([]byte, error)
	GetMacrosList(ctx context.Context) ([]byte, error)
	GetNDISources(ctx context.Context) ([]byte, error)
	GetDDRPlaylist(ctx context.Context) ([]byte, error)
	GetDDRTimecode(ctx context.Context) ([]byte, error)

	// DataLink()
	// File()
	VideoPreview(ctx context.Context, name string, xres, yres, quality int) (image.Image, error)
	// Icon()

	// Dictonary stuff
	// Dictionary(key string) error
	MetaData(ctx context.Context) (*Metadata, error)
	// Tally()
	// Switcher()
	// Buffer()
	// SwitcherUIEffect()
	ShortcutStates(ctx context.Context) (*ShortcutStates, error)
	// DDRPlaylist
	// DDRTimecode
	// AudioMixer
	// AudioBins
	// FileBrowser
	// MacrosList
	// NDISources

	// Websocket stuff
	// v1/audio_notifications
	// v1/change_notifications
	// v1/shortcut_notifications
	// v1/shortcut_state_notifications
	// v1/timecode_notifications
	// v1/vu_notifications
	VUNotifications(ctx context.Context, name string, cb func(msg string)) error
	ChangeNotifications(ctx context.Context, cb func(msg string)) error
	TimecodeNotifications(ctx context.Context, cb func(msg string)) error
	ShortcutNotifications(ctx context.Context, cb func(msg string)) error

	// Additional WebSocket features
	AudioNotifications(ctx context.Context, name string, cb func(data []byte)) error
	ShortcutStateNotifications(ctx context.Context, cb func(msg string)) error
	VideoNotifications(ctx context.Context, name string, xres, yres, quality int, cb func(data []byte)) error

	// VideoPreview(name string,xres int,yres int,quality string)
	// http://systemnameoripaddress/v1/video_notifications?name=NAME&xres=RESX&yres=RESY&q=QUALITY

	// Audio Send/Receive
	// AudioOutput()
	// AudioAUX()
	// AudioPhones()
	// http://systemnameoripaddress/v1/audio_notifications?name=NAME

	// VU meter
	// VUMeter()
	// http://systemnameoripaddress/
}

type clientV1 struct {
	host     *url.URL
	user     string
	password string

	// Websocket connection cache
	wsConnections map[string]*websocket.Conn
	wsMutex       sync.RWMutex
	wsCallbacks   map[string][]func(msg string)
	cbMutex       sync.RWMutex
}

func NewClientV1(host, user, password string) (ClientV1, error) {
	u, err := url.Parse(fmt.Sprintf("http://%s/v1", host))
	if err != nil {
		return nil, err
	}

	return &clientV1{
		host:          u,
		user:          user,
		password:      password,
		wsConnections: make(map[string]*websocket.Conn),
		wsCallbacks:   make(map[string][]func(msg string)),
	}, nil
}

func (c *clientV1) get(ctx context.Context, endpoint string, queries map[string]string) ([]byte, error) {
	u := *c.host
	u.Path = path.Join(u.Path, endpoint)
	q := u.Query()
	for k, v := range queries {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	uri := u.String()

	client := &http.Client{
		Transport: &digest.Transport{
			Username: c.user,
			Password: c.password,
		},
	}
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *clientV1) GetLive(ctx context.Context) bool {
	b, err := c.get(ctx, "./live", nil)
	if err != nil {
		return false
	}

	switch string(b) {
	case "TRUE":
		return true
	case "FALSE":
		return false
	default:
		return false // Unknown case
	}
}

func (c *clientV1) GetProduct(ctx context.Context) (info *ProductInformation, err error) {
	b, err := c.get(ctx, "./version", nil)
	if err != nil {
		return nil, err
	}

	ret := &ProductInformation{}
	if err := xml.Unmarshal(b, ret); err != nil {
		fmt.Println("body:", b)
		return nil, err
	}
	return ret, nil
}

func (c *clientV1) ShortcutHTTP(ctx context.Context, name string, kv map[string]string) error {
	if kv == nil {
		kv = make(map[string]string)
	}
	kv["name"] = name

	_, err := c.get(ctx, "./shortcut", kv)
	if err != nil {
		return err
	}
	return nil
}
func (c *clientV1) ShortcutWS(ctx context.Context, name string, kv map[string]string) error {
	// WebSocket経由でのショートカット実行
	// v1/shortcut_state を使用してショートカットを送信
	conn, err := c.getWebSocketConnection(ctx, "shortcut_state")
	if err != nil {
		return err
	}

	// XML形式でショートカットを構築（encoding/xmlを使用）
	type ShortcutWithAttrs struct {
		XMLName xml.Name `xml:"shortcut"`
		Name    string   `xml:"name,attr"`
		Value   string   `xml:"value,attr,omitempty"`
		// 動的属性を扱うためのカスタムマーシャリング
		ExtraAttrs map[string]string `xml:"-"`
	}

	// MarshalXMLメソッドを実装
	shortcut := ShortcutWithAttrs{
		Name:       name,
		Value:      "",
		ExtraAttrs: kv,
	}

	// XMLを生成（動的属性を含む）
	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)

	// 開始タグを作成
	start := xml.StartElement{
		Name: xml.Name{Local: "shortcut"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "name"}, Value: shortcut.Name},
		},
	}

	if shortcut.Value != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "value"}, Value: shortcut.Value})
	}

	// 追加の属性を追加
	for k, v := range shortcut.ExtraAttrs {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: k}, Value: v})
	}

	// XMLをエンコード
	if err := encoder.EncodeToken(start); err != nil {
		return fmt.Errorf("XML encoding failed: %w", err)
	}
	if err := encoder.EncodeToken(start.End()); err != nil {
		return fmt.Errorf("XML encoding failed: %w", err)
	}
	if err := encoder.Flush(); err != nil {
		return fmt.Errorf("XML encoding failed: %w", err)
	}

	// WebSocket経由でショートカットを送信
	return conn.Write(ctx, websocket.MessageText, buf.Bytes())
}
func (c *clientV1) Trigger(ctx context.Context, name string) error {
	// トリガー実行 - /v1/trigger?name=NAME を使用
	_, err := c.get(ctx, "./trigger", map[string]string{"name": name})
	return err
}

func (c *clientV1) VideoPreview(ctx context.Context, name string, xres, yres, quality int) (image.Image, error) {
	b, err := c.get(ctx, "./image", map[string]string{
		"name": name,
		"xres": strconv.Itoa(xres),
		"yres": strconv.Itoa(yres),
		"q":    strconv.Itoa(quality),
	})
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(b)
	i, err := jpeg.Decode(r)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (c *clientV1) ShortcutStates(ctx context.Context) (*ShortcutStates, error) {
	mp := make(map[string]string)
	mp["name"] = "shortcut_states"

	ret := &ShortcutStates{}

	b, err := c.get(ctx, "./dictionary", mp)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(b, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *clientV1) MetaData(ctx context.Context) (*Metadata, error) {
	mp := make(map[string]string)
	mp["key"] = "metadata"

	ret := &Metadata{}

	b, err := c.get(ctx, "./dictionary", mp)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(b, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *clientV1) ChangeNotifications(ctx context.Context, cb func(msg string)) error {
	conn, err := c.getOrCreateWebSocketConnection(ctx, "change_notifications", nil)
	if err != nil {
		return err
	}

	// Add callback for this notification type
	c.addCallback("change_notifications", nil, cb)

	// Start message handler if not already started
	c.startMessageHandler(ctx, conn, "change_notifications", nil)

	return nil
}

func (c *clientV1) TimecodeNotifications(ctx context.Context, cb func(msg string)) error {
	conn, err := c.getOrCreateWebSocketConnection(ctx, "timecode_notifications", nil)
	if err != nil {
		return err
	}

	// Add callback for this notification type
	c.addCallback("timecode_notifications", nil, cb)

	// Start message handler if not already started
	c.startMessageHandler(ctx, conn, "timecode_notifications", nil)

	return nil
}

func (c *clientV1) ShortcutNotifications(ctx context.Context, cb func(msg string)) error {
	conn, err := c.getOrCreateWebSocketConnection(ctx, "shortcut_notifications", nil)
	if err != nil {
		return err
	}

	// Add callback for this notification type
	c.addCallback("shortcut_notifications", nil, cb)

	// Start message handler if not already started
	c.startMessageHandler(ctx, conn, "shortcut_notifications", nil)

	return nil
}

func (c *clientV1) VUNotifications(ctx context.Context, name string, cb func(msg string)) error {
	params := map[string]string{"name": name}
	conn, err := c.getOrCreateWebSocketConnection(ctx, "vu_notifications", params)
	if err != nil {
		return err
	}

	// Add callback for this notification type
	c.addCallback("vu_notifications", params, cb)

	// Start message handler if not already started
	c.startMessageHandler(ctx, conn, "vu_notifications", params)

	return nil
}

// getOrCreateWebSocketConnection gets an existing WebSocket connection or creates a new one
func (c *clientV1) getOrCreateWebSocketConnection(ctx context.Context, notificationType string, params map[string]string) (*websocket.Conn, error) {
	// Create a unique key for this connection type and parameters
	key := notificationType
	for k, v := range params {
		key += ":" + k + "=" + v
	}

	c.wsMutex.RLock()
	conn, exists := c.wsConnections[key]
	c.wsMutex.RUnlock()

	if exists && conn != nil {
		return conn, nil
	}

	// Create new connection
	uri, cred, err := c.getCred(ctx, notificationType)
	if err != nil {
		return nil, err
	}

	// Add query parameters if provided
	if len(params) > 0 {
		u, err := url.Parse(uri)
		if err != nil {
			return nil, err
		}
		q := u.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
		uri = u.String()
	}

	conn, _, err = websocket.Dial(ctx, uri, &websocket.DialOptions{
		HTTPHeader: cred,
	})
	if err != nil {
		return nil, err
	}

	// Store the connection
	c.wsMutex.Lock()
	c.wsConnections[key] = conn
	c.wsMutex.Unlock()

	// Start ping goroutine for this connection
	go c.startPingRoutine(ctx, conn)

	return conn, nil
}

// addCallback adds a callback function for a specific notification type
func (c *clientV1) addCallback(notificationType string, params map[string]string, callback func(msg string)) {
	key := notificationType
	for k, v := range params {
		key += ":" + k + "=" + v
	}

	c.cbMutex.Lock()
	defer c.cbMutex.Unlock()
	c.wsCallbacks[key] = append(c.wsCallbacks[key], callback)
}

// startMessageHandler starts a goroutine to handle incoming messages for a connection
func (c *clientV1) startMessageHandler(ctx context.Context, conn *websocket.Conn, notificationType string, params map[string]string) {
	key := notificationType
	for k, v := range params {
		key += ":" + k + "=" + v
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := conn.Read(ctx)
				if err != nil {
					log.Printf("WebSocket read error for %s: %v", key, err)
					return
				}

				// Call all registered callbacks for this connection
				c.cbMutex.RLock()
				callbacks := c.wsCallbacks[key]
				c.cbMutex.RUnlock()

				for _, callback := range callbacks {
					callback(string(message))
				}
			}
		}
	}()
}

// startPingRoutine starts a ping routine to keep the connection alive
func (c *clientV1) startPingRoutine(ctx context.Context, conn *websocket.Conn) {
	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := conn.Write(ctx, websocket.MessageText, []byte("ping")); err != nil {
				log.Printf("WebSocket ping error: %v", err)
				return
			}
		}
	}
}

func (c *clientV1) getCred(ctx context.Context, path string) (string, http.Header, error) {
	uri := c.host.String()
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return "", nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer res.Body.Close()
	if _, err := io.Copy(io.Discard, res.Body); err != nil {
		return "", nil, err
	}

	header := res.Header.Get("WWW-Authenticate")
	chal, _ := digest.ParseChallenge(header)

	// use it to create credentials for the next request
	cred, _ := digest.Digest(chal, digest.Options{
		Username: c.user,
		Password: c.password,
		Method:   res.Request.Method,
		URI:      uri,
		Count:    2,
	})
	uri = fmt.Sprintf("ws://%s/v1/%s", c.host.Host, path)
	res.Request.RequestURI = uri
	res.Request.Header.Set("Authorization", cred.String())

	return uri, res.Request.Header, nil
}

// Basic Switcher Controls
func (c *clientV1) Take(ctx context.Context) error {
	return c.ShortcutHTTP(ctx, "take", nil)
}

func (c *clientV1) AutoTransition(ctx context.Context) error {
	return c.ShortcutHTTP(ctx, "auto_transition", nil)
}

func (c *clientV1) SetSourceToPreview(ctx context.Context, source string) error {
	return c.ShortcutHTTP(ctx, "set_preview", map[string]string{"source": source})
}

func (c *clientV1) SetSourceToProgram(ctx context.Context, source string) error {
	return c.ShortcutHTTP(ctx, "set_program", map[string]string{"source": source})
}

func (c *clientV1) SetSourceToME(ctx context.Context, me, source string) error {
	return c.ShortcutHTTP(ctx, "set_me", map[string]string{"me": me, "source": source})
}

func (c *clientV1) SetDSKOnAir(ctx context.Context, dsk string, onAir bool) error {
	value := "false"
	if onAir {
		value = "true"
	}
	return c.ShortcutHTTP(ctx, "set_dsk_onair", map[string]string{"dsk": dsk, "onair": value})
}

func (c *clientV1) SetTransitionSelection(ctx context.Context, transition string) error {
	return c.ShortcutHTTP(ctx, "set_transition", map[string]string{"transition": transition})
}

// Media Controls
func (c *clientV1) MediaPlay(ctx context.Context, media string) error {
	return c.ShortcutHTTP(ctx, "media_play", map[string]string{"media": media})
}

func (c *clientV1) MediaStop(ctx context.Context, media string) error {
	return c.ShortcutHTTP(ctx, "media_stop", map[string]string{"media": media})
}

func (c *clientV1) MediaRewind(ctx context.Context, media string) error {
	return c.ShortcutHTTP(ctx, "media_rewind", map[string]string{"media": media})
}

func (c *clientV1) MediaFastForward(ctx context.Context, media string) error {
	return c.ShortcutHTTP(ctx, "media_fastforward", map[string]string{"media": media})
}

func (c *clientV1) MediaPause(ctx context.Context, media string) error {
	return c.ShortcutHTTP(ctx, "media_pause", map[string]string{"media": media})
}

// Recording and Streaming
func (c *clientV1) RecordToggle(ctx context.Context) error {
	return c.ShortcutHTTP(ctx, "record", nil)
}

func (c *clientV1) RecordStart(ctx context.Context) error {
	return c.ShortcutHTTP(ctx, "start-record", nil)
}

func (c *clientV1) RecordStop(ctx context.Context) error {
	return c.ShortcutHTTP(ctx, "stop-record", nil)
}

func (c *clientV1) StreamToggle(ctx context.Context) error {
	return c.ShortcutHTTP(ctx, "stream", nil)
}

func (c *clientV1) StreamStart(ctx context.Context) error {
	return c.ShortcutHTTP(ctx, "start-stream", nil)
}

func (c *clientV1) StreamStop(ctx context.Context) error {
	return c.ShortcutHTTP(ctx, "stop-stream", nil)
}

// Macros
func (c *clientV1) RunSystemMacro(ctx context.Context, macro string) error {
	return c.ShortcutHTTP(ctx, "system_macro", map[string]string{"macro": macro})
}

func (c *clientV1) RunCustomMacro(ctx context.Context, macro string) error {
	return c.ShortcutHTTP(ctx, "custom_macro", map[string]string{"macro": macro})
}

// Advanced Macro Controls
func (c *clientV1) PlayMacroByName(ctx context.Context, name string) error {
	return c.ShortcutHTTP(ctx, "play_macro_byname", map[string]string{"name": name})
}

func (c *clientV1) PlayMacroByID(ctx context.Context, id string) error {
	return c.ShortcutHTTP(ctx, "play_macro_byid", map[string]string{"id": id})
}

func (c *clientV1) RecordMacroByName(ctx context.Context, name string) error {
	return c.ShortcutHTTP(ctx, "record_macro_byname", map[string]string{"name": name})
}

func (c *clientV1) StopRecordMacro(ctx context.Context) error {
	return c.ShortcutHTTP(ctx, "stop_record_macro", nil)
}

func (c *clientV1) StopMacroByName(ctx context.Context, name string) error {
	return c.ShortcutHTTP(ctx, "stop_macro_byname", map[string]string{"name": name})
}

func (c *clientV1) StopMacroByID(ctx context.Context, id string) error {
	return c.ShortcutHTTP(ctx, "stop_macro_byid", map[string]string{"id": id})
}

func (c *clientV1) MacroIsEnabledByName(ctx context.Context, name string) (bool, error) {
	b, err := c.get(ctx, "./macro_enabled", map[string]string{"name": name})
	if err != nil {
		return false, err
	}

	switch string(b) {
	case "TRUE", "true", "1":
		return true, nil
	case "FALSE", "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("unknown macro enabled status: %s", string(b))
	}
}

func (c *clientV1) MacroIsEnabledByID(ctx context.Context, id string) (bool, error) {
	b, err := c.get(ctx, "./macro_enabled", map[string]string{"id": id})
	if err != nil {
		return false, err
	}

	switch string(b) {
	case "TRUE", "true", "1":
		return true, nil
	case "FALSE", "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("unknown macro enabled status: %s", string(b))
	}
}

func (c *clientV1) GetMacroIDByName(ctx context.Context, name string) (string, error) {
	b, err := c.get(ctx, "./get_macroid", map[string]string{"name": name})
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Output Control
func (c *clientV1) SetMixOutput(ctx context.Context, output string) error {
	return c.ShortcutHTTP(ctx, "set_mix_output", map[string]string{"output": output})
}

// DataLink
func (c *clientV1) SetDataLinkValue(ctx context.Context, key, value string) error {
	return c.ShortcutHTTP(ctx, "set_datalink", map[string]string{"key": key, "value": value})
}

func (c *clientV1) GetDataLinkValue(ctx context.Context, key string) (string, error) {
	b, err := c.get(ctx, "./datalink", map[string]string{"key": key})
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Feedback and Status
func (c *clientV1) GetTally(ctx context.Context) (*Tally, error) {
	b, err := c.get(ctx, "./tally", nil)
	if err != nil {
		return nil, err
	}

	ret := &Tally{}
	if err := xml.Unmarshal(b, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *clientV1) GetRecordingStatus(ctx context.Context) (bool, error) {
	b, err := c.get(ctx, "./recording_status", nil)
	if err != nil {
		return false, err
	}

	switch string(b) {
	case "TRUE", "true", "1":
		return true, nil
	case "FALSE", "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("unknown recording status: %s", string(b))
	}
}

func (c *clientV1) GetStreamingStatus(ctx context.Context) (bool, error) {
	b, err := c.get(ctx, "./streaming_status", nil)
	if err != nil {
		return false, err
	}

	switch string(b) {
	case "TRUE", "true", "1":
		return true, nil
	case "FALSE", "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("unknown streaming status: %s", string(b))
	}
}

func (c *clientV1) GetMediaStatus(ctx context.Context, media string) (string, error) {
	b, err := c.get(ctx, "./media_status", map[string]string{"media": media})
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *clientV1) GetSourceStatus(ctx context.Context) (program, preview string, err error) {
	// プログラムソースの取得
	progBytes, err := c.get(ctx, "./program_source", nil)
	if err != nil {
		return "", "", err
	}
	program = string(progBytes)

	// プレビューソースの取得
	prevBytes, err := c.get(ctx, "./preview_source", nil)
	if err != nil {
		return "", "", err
	}
	preview = string(prevBytes)

	return program, preview, nil
}

// XMLエスケープ関数
func xmlEscape(s string) string {
	return html.EscapeString(s)
}

// Dictionary Queries
func (c *clientV1) GetDictionary(ctx context.Context, key string) ([]byte, error) {
	return c.get(ctx, "./dictionary", map[string]string{"key": key})
}

func (c *clientV1) GetSwitcherInfo(ctx context.Context) ([]byte, error) {
	return c.GetDictionary(ctx, "switcher")
}

func (c *clientV1) GetBufferInfo(ctx context.Context) ([]byte, error) {
	return c.GetDictionary(ctx, "buffer")
}

func (c *clientV1) GetSwitcherUIEffects(ctx context.Context) ([]byte, error) {
	return c.GetDictionary(ctx, "switcher_ui_effects")
}

func (c *clientV1) GetAudioMixerInfo(ctx context.Context) ([]byte, error) {
	return c.GetDictionary(ctx, "audiomixer")
}

func (c *clientV1) GetAudioBinsInfo(ctx context.Context) ([]byte, error) {
	return c.GetDictionary(ctx, "audio_bins")
}

func (c *clientV1) GetFileBrowserInfo(ctx context.Context) ([]byte, error) {
	return c.GetDictionary(ctx, "filebrowser")
}

func (c *clientV1) GetMacrosList(ctx context.Context) ([]byte, error) {
	return c.GetDictionary(ctx, "macros_list")
}

func (c *clientV1) GetNDISources(ctx context.Context) ([]byte, error) {
	return c.GetDictionary(ctx, "ndi_sources")
}

func (c *clientV1) GetDDRPlaylist(ctx context.Context) ([]byte, error) {
	return c.GetDictionary(ctx, "ddr_playlist")
}

func (c *clientV1) GetDDRTimecode(ctx context.Context) ([]byte, error) {
	return c.GetDictionary(ctx, "ddr_timecode")
}

// WebSocket connection management
func (c *clientV1) getWebSocketConnection(ctx context.Context, path string) (*websocket.Conn, error) {
	c.wsMutex.Lock()
	defer c.wsMutex.Unlock()

	// 既存の接続があるかチェック
	if conn, exists := c.wsConnections[path]; exists {
		return conn, nil
	}

	// 新しい接続を作成
	uri, cred, err := c.getCred(ctx, path)
	if err != nil {
		return nil, err
	}

	conn, resp, err := websocket.Dial(ctx, uri, &websocket.DialOptions{
		HTTPHeader: cred,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return nil, err
	}

	// 接続をキャッシュに保存
	if c.wsConnections == nil {
		c.wsConnections = make(map[string]*websocket.Conn)
	}
	c.wsConnections[path] = conn

	// 接続維持のためのpingルーチンを開始
	go c.keepAlive(ctx, conn)

	return conn, nil
}

func (c *clientV1) keepAlive(ctx context.Context, conn *websocket.Conn) {
	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := conn.Ping(ctx); err != nil {
				log.Printf("WebSocket ping failed: %v", err)
				return
			}
		}
	}
}

// Additional WebSocket features
func (c *clientV1) AudioNotifications(ctx context.Context, name string, cb func(data []byte)) error {
	uri, cred, err := c.getCred(ctx, "audio_notifications")
	if err != nil {
		return err
	}
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Add("name", name)
	u.RawQuery = q.Encode()

	conn, resp, err := websocket.Dial(ctx, u.String(), &websocket.DialOptions{
		HTTPHeader: cred,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := conn.Read(ctx)
				if err != nil {
					log.Println("audio read:", err)
					return
				}
				cb(message)
			}
		}
	}()

	go c.keepAlive(ctx, conn)
	return nil
}

func (c *clientV1) ShortcutStateNotifications(ctx context.Context, cb func(msg string)) error {
	uri, cred, err := c.getCred(ctx, "shortcut_state_notifications")
	if err != nil {
		return err
	}

	conn, resp, err := websocket.Dial(ctx, uri, &websocket.DialOptions{
		HTTPHeader: cred,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := conn.Read(ctx)
				if err != nil {
					log.Println("shortcut state read:", err)
					return
				}
				cb(string(message))
			}
		}
	}()

	go c.keepAlive(ctx, conn)
	return nil
}

func (c *clientV1) VideoNotifications(ctx context.Context, name string, xres, yres, quality int, cb func(data []byte)) error {
	uri, cred, err := c.getCred(ctx, "video_notifications")
	if err != nil {
		return err
	}
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Add("name", name)
	q.Add("xres", strconv.Itoa(xres))
	q.Add("yres", strconv.Itoa(yres))
	q.Add("q", strconv.Itoa(quality))
	u.RawQuery = q.Encode()

	conn, resp, err := websocket.Dial(ctx, u.String(), &websocket.DialOptions{
		HTTPHeader: cred,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := conn.Read(ctx)
				if err != nil {
					log.Println("video read:", err)
					return
				}
				cb(message)
			}
		}
	}()

	go c.keepAlive(ctx, conn)
	return nil
}
