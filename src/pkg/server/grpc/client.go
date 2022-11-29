package grpc

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Client struct {
	endPoint string
}

func NewClient(e string) *Client {

	return &Client{
		endPoint: e,
	}
}

func (c *Client) Dial() (*grpc.ClientConn, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, c.endPoint)
	if err != nil {
		return nil, errors.Wrap(err, "grpc dial error")
	}
	return conn, nil
}
