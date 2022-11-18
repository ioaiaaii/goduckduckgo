package store

import (
	"context"
	"goduckduckgo/pkg/store/storepb"

	db "github.com/go-pg/pg"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DBStore struct {
	client *db.DB
}

func NewDBStore(c *db.DB) (*DBStore, error) {

	dbs := &DBStore{
		client: c,
	}
	return dbs, nil
}

// CreateTodo creates a todo given a description
func (s *DBStore) CreateTodo(ctx context.Context, req *storepb.CreateTodoRequest) (*storepb.CreateTodoResponse, error) {
	req.Item.Id = uuid.NewV4().String()
	err := s.client.Insert(req.Item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not insert item into the database: %s", err)
	}

	return &storepb.CreateTodoResponse{Id: req.Item.Id}, nil
}
