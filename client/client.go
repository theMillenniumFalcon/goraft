package client

import (
	"context"
	"net"

	"github.com/theMillenniumFalcon/goraft/proto"
)

type Options struct{}

type Client struct {
	conn net.Conn
}

func New(endpoint string, opts Options) (*Client, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Set(ctx context.Context, key []byte, value []byte, timeToLive int) (any, error) {
	cmd := &proto.CommandSet{
		Key:        key,
		Value:      value,
		TimeToLive: timeToLive,
	}

	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
