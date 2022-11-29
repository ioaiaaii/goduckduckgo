package store

import (
	"context"
	"goduckduckgo/pkg/store/storepb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Target string
	Conn   *grpc.ClientConn
	Client storepb.StoreClient
}

func RegisterStoreClient(target string) *Client {
	return &Client{
		Target: target,
	}
}

func (s *Client) Dial(target string) (*grpc.ClientConn, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	s.Conn = conn

	return conn, nil
}

// func (s *StoreClient) NewStoreClient(target string, g *grpc.ClientConn) (storepb.StoreClient, error) {
// 	c := storepb.NewStoreClient(g)
// 	return c, nil
// }

func (s *Client) NewStoreClient() (storepb.StoreClient, error) {

	var (
		conn *grpc.ClientConn
		err  error
	)

	conn, err = s.Dial(s.Target)
	if err != nil {
		return nil, err
	}

	c := storepb.NewStoreClient(conn)
	s.Client = c
	return c, nil

}
