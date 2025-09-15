package newtek

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
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
	panic("not implemented")
}
func (c *clientV1) Trigger(ctx context.Context, name string) error {
	panic("not implemented")
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
