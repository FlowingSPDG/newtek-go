package newtek

import (
	"bytes"
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
	"time"

	"github.com/gorilla/websocket"
	"github.com/icholy/digest"
)

type ClientV1 interface {
	GetLive() bool
	GetProduct() (info *ProductInformation, err error)
	ShortcutHTTP(name string, kv map[string]string) error
	ShortcutWS(name string, kv map[string]string) error
	Trigger(name string) error
	// DataLink()
	// File()
	VideoPreview(name string, xres, yres, quality int) (image.Image, error)
	// Icon()

	// Dictonary stuff
	// Dictionary(key string) error
	MetaData() (*Metadata, error)
	// Tally()
	// Switcher()
	// Buffer()
	// SwitcherUIEffect()
	ShortcutStates() (*ShortcutStates, error)
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
	VUNotifications(name string, cb func(msg string)) error
	ChangeNotifications(cb func(msg string)) error
	TimecodeNotifications(cb func(msg string)) error
	ShortcutNotifications(cb func(msg string)) error

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

	// Websocket notifications
	// audioNotifications        *websocket.Conn
	changeNotifications *websocket.Conn
	// shortcutNotifications     *websocket.Conn
	// shortcutSafeNotifications *websocket.Conn
	// timecodeNotifications     *websocket.Conn
	// vuNotifications *websocket.Conn
}

func NewClientV1(host, user, password string) (ClientV1, error) {
	u, err := url.Parse(fmt.Sprintf("http://%s/v1", host))
	if err != nil {
		return nil, err
	}

	return &clientV1{
		host:     u,
		user:     user,
		password: password,
	}, nil
}

func (c *clientV1) get(endpoint string, queries map[string]string) ([]byte, error) {
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
	resp, err := client.Get(uri)
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

func (c *clientV1) GetLive() bool {
	b, err := c.get("./live", nil)
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

func (c *clientV1) GetProduct() (info *ProductInformation, err error) {
	b, err := c.get("./version", nil)
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

func (c *clientV1) ShortcutHTTP(name string, kv map[string]string) error {
	if kv == nil {
		kv = make(map[string]string)
	}
	kv["name"] = name

	_, err := c.get("./shortcut", kv)
	if err != nil {
		return err
	}
	return nil
}
func (c *clientV1) ShortcutWS(name string, kv map[string]string) error {
	panic("not implemented")
}
func (c *clientV1) Trigger(name string) error {
	panic("not implemented")
}

func (c *clientV1) VideoPreview(name string, xres, yres, quality int) (image.Image, error) {
	b, err := c.get("./image", map[string]string{
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

func (c *clientV1) ShortcutStates() (*ShortcutStates, error) {
	mp := make(map[string]string)
	mp["name"] = "shortcut_states"

	ret := &ShortcutStates{}

	b, err := c.get("./dictionary", mp)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(b, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *clientV1) MetaData() (*Metadata, error) {
	mp := make(map[string]string)
	mp["key"] = "metadata"

	ret := &Metadata{}

	b, err := c.get("./dictionary", mp)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(b, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *clientV1) ChangeNotifications(cb func(msg string)) error {
	uri, cred, err := c.getCred("change_notifications")
	if err != nil {
		return err
	}

	d := websocket.DefaultDialer
	con, resp, err := d.Dial(uri, cred)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return err
	}
	c.changeNotifications = con

	// TODO: Support context
	go func() {
		for {
			_, message, err := con.ReadMessage()
			if err != nil {
				log.Println("read:", err)
			}
			cb(string(message))
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(15))
			con.WriteMessage(websocket.PingMessage, []byte("ping"))
		}
	}()

	return nil
}

func (c *clientV1) TimecodeNotifications(cb func(msg string)) error {
	uri, cred, err := c.getCred("timecode_notifications")
	if err != nil {
		return err
	}
	d := websocket.DefaultDialer
	con, resp, err := d.Dial(uri, cred)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return err
	}
	c.changeNotifications = con

	// TODO: Support context
	go func() {
		for {
			_, message, err := con.ReadMessage()
			if err != nil {
				log.Println("read:", err)
			}
			cb(string(message))
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(15))
			con.WriteMessage(websocket.PingMessage, []byte("ping"))
		}
	}()

	return nil
}

func (c *clientV1) ShortcutNotifications(cb func(msg string)) error {
	uri, cred, err := c.getCred("shortcut_notifications")
	if err != nil {
		return err
	}
	d := websocket.DefaultDialer
	con, resp, err := d.Dial(uri, cred)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return err
	}
	c.changeNotifications = con

	// TODO: Support context
	go func() {
		for {
			_, message, err := con.ReadMessage()
			if err != nil {
				log.Println("read:", err)
			}
			cb(string(message))
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(15))
			con.WriteMessage(websocket.PingMessage, []byte("ping"))
		}
	}()

	return nil
}

func (c *clientV1) VUNotifications(name string, cb func(msg string)) error {
	uri, cred, err := c.getCred("vu_notifications")
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
	d := websocket.DefaultDialer
	con, resp, err := d.Dial(uri, cred)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return err
	}
	c.changeNotifications = con

	// TODO: Support context
	go func() {
		for {
			_, message, err := con.ReadMessage()
			if err != nil {
				log.Println("read:", err)
			}
			cb(string(message))
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(15))
			con.WriteMessage(websocket.PingMessage, []byte("ping"))
		}
	}()

	return nil
}

func (c *clientV1) getCred(path string) (string, http.Header, error) {
	uri := c.host.String()
	res, err := http.Get(uri)
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
