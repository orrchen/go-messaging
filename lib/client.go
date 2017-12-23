package lib

import (
	"net"
	"bufio"
	"io"
	"log"
)
type ConnectionEventType string

const
(
	CONNECTION_EVENT_TYPE_NEW_CONNECTION        ConnectionEventType = "new_connection"
	CONNECTION_EVENT_TYPE_CONNECTION_TERMINATED ConnectionEventType = "connection_terminated"
	CONNECTION_EVENT_TYPE_CONNECTION_GENERAL_ERROR ConnectionEventType = "general_error"
)

// Client holds info about connection
type Client struct {
	Uid  string /* client is responsible of generating a uinque uid for each request, it will be sent in the response from the server so that client will know what request generated this response */
	DeviceUid string /* a unique id generated from the client itself */
	conn net.Conn
	onConnectionEvent func(c *Client, eventType ConnectionEventType, e error) /* function for handling new connections */
	onDataEvent func(c *Client, data []byte) /* function for handling new date events */
}

func NewClient(conn net.Conn,onConnectionEvent func(c *Client,eventType ConnectionEventType, e error), onDataEvent func(c *Client, data []byte)) *Client {
	return  &Client{
		conn: conn,
		onConnectionEvent: onConnectionEvent,
		onDataEvent: onDataEvent,
	}
}

// Read client data from channel
func (c *Client) listen() {
	reader := bufio.NewReader(c.conn)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)

		switch err {
		case io.EOF:
			// connection terminated
			c.conn.Close()
			c.onConnectionEvent(c,CONNECTION_EVENT_TYPE_CONNECTION_TERMINATED, err)
			return
		case nil:
			// new data available
			c.onDataEvent(c, buf[:n])
		default:
			log.Fatalf("Receive data failed:%s", err)
			c.conn.Close()
			c.onConnectionEvent(c, CONNECTION_EVENT_TYPE_CONNECTION_GENERAL_ERROR, err)
			return
		}
	}
}

// Send text message to client
func (c *Client) Send(message []byte) error {
	_, err := c.conn.Write(message)
	return err
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	if c.conn!=nil {
		return c.conn.Close()
	}
	return nil

}
