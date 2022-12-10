package newtek

import (
	"encoding/xml"
	"net"
)

type clientV1TCP struct {
	conn net.Conn // TCP Connection
}

// Register implements ClientV1TCP
func (c *clientV1TCP) Register(name string) error {
	r := register{
		Name: name, // "NTK_states"?
	}
	b, err := xml.Marshal(r)
	if err != nil {
		return err
	}
	return c.SendBytes(b)
}

// Unregister implements ClientV1TCP
func (c *clientV1TCP) Unregister(name string) error {
	r := unregister{
		Name: name, // "NTK_states"?
	}
	b, err := xml.Marshal(r)
	if err != nil {
		return err
	}
	return c.SendBytes(b)
}

// Send implements ClientV1TCP
func (c *clientV1TCP) SendBytes(data []byte) error {
	data = append(data, "\n"...)
	_, err := c.conn.Write(data)
	return err
}

// Send implements ClientV1TCP
func (c *clientV1TCP) Send(data string) error {
	return c.SendBytes([]byte(data))
}

func (c *clientV1TCP) Close() error {
	return c.conn.Close()
}

// Shortcut implements ClientV1TCP
func (c *clientV1TCP) Shortcut(s ShortcutStates) error {
	panic("unimplemented")
}

type register struct {
	XMLName xml.Name `xml:"register"`
	Text    string   `xml:",chardata"`
	Name    string   `xml:"name"`
}

type unregister struct {
	XMLName xml.Name `xml:"unregister"`
	Text    string   `xml:",chardata"`
	Name    string   `xml:"name"`
}

type ClientV1TCP interface {
	Close() error                 // Close Connection
	Register(name string) error   // send NTK_states
	Unregister(name string) error // send NTK_states
	Send(data string) error
	SendBytes(data []byte) error
	Shortcut(s ShortcutStates) error // should include shortcut XML or name and K/V

	// TODO: Tally from reading packet
}

func NewClientV1TCP(host string) (ClientV1TCP, error) {
	p := net.JoinHostPort(host, "5951")
	c, err := net.Dial("tcp", p)
	if err != nil {
		return nil, err
	}

	// go io.Copy(os.Stdout, c)
	ret := &clientV1TCP{
		conn: c,
	}

	return ret, nil
}
