package store

import (
	"context"
	"encoding/json"
	"goduckduckgo/internal/storage/db"
	"goduckduckgo/pkg/duckduckgo/typespb"
	"goduckduckgo/pkg/store/storepb"
	"net/http"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	err error
)

type Server struct {
	client *db.Database
	storepb.UnimplementedStoreServer
}

// Returns a pointer to the newly created Store
func NewStore(c *db.Database) *Server {

	return &Server{
		client: c,
	}

}

func (s *Server) Create(_ context.Context, req *storepb.CreateRequest) (*storepb.CreateResponse, error) {

	ja, _ := json.Marshal(req.Answer)
	q := db.DDGQueryTable{
		Query:  req.Query,
		Answer: ja,
	}

	err = s.client.GormClient.Create(&q).Error

	if err != nil {
		return &storepb.CreateResponse{
				Status: http.StatusInternalServerError,
			},
			status.Errorf(codes.Internal, "Could not insert item into the database: %s", err)
	}

	return &storepb.CreateResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Read(_ context.Context, req *storepb.ReadRequest) (*storepb.ReadResponse, error) {

	var (
		resultT db.DDGQueryTable
		a       *typespb.QueryPayload
	)

	err = s.client.GormClient.Where("query = ?", req.Query).Find(&resultT).Error

	if err != nil {
		return &storepb.ReadResponse{Status: http.StatusInternalServerError, Error: err.Error()}, status.Errorf(codes.Internal, "Error during read: %s", err)
	}

	err = json.Unmarshal(resultT.Answer, &a)
	log.Info().Msgf("New Raw %v", resultT)

	if err != nil {
		return &storepb.ReadResponse{Status: http.StatusInternalServerError, Error: err.Error()}, status.Errorf(codes.Internal, "Error during marshalling DB Response: %s", err)
	}

	return &storepb.ReadResponse{
		Status: http.StatusOK,
		Answer: a,
	}, nil
}

func (s *Server) Delete(_ context.Context, req *storepb.DeleteRequest) (*storepb.DeleteResponse, error) {

	var (
		deleteT db.DDGQueryTable
	)
	err = s.client.GormClient.Where("query = ?", req.Query).Find(&deleteT).Error
	if err != nil {
		return &storepb.DeleteResponse{Status: http.StatusInternalServerError, Error: err.Error()}, status.Errorf(codes.Internal, "Error during searching query to delete: %s", err)
	}

	if deleteT.Answer == nil {
		log.Warn().Msgf("Query %s not found, bypassing deletion.", req.Query)
		return &storepb.DeleteResponse{
			Status: http.StatusOK,
		}, nil
	}

	err = s.client.GormClient.Delete(&deleteT).Error
	if err != nil {
		return &storepb.DeleteResponse{Status: http.StatusInternalServerError, Error: err.Error()}, status.Errorf(codes.Internal, "Error during query deletion:: %s", err)
	}

	return &storepb.DeleteResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) Update(_ context.Context, req *storepb.UpdateRequest) (*storepb.UpdateResponse, error) {

	var (
		updateT db.DDGQueryTable
	)

	ura, _ := json.Marshal(req.Answer)
	updateQuery := db.DDGQueryTable{
		Query:  req.Query,
		Answer: ura,
	}

	err = s.client.GormClient.Model(updateT).Where("query = ?", req.Query).Updates(&updateQuery).Error

	if err != nil {
		return &storepb.UpdateResponse{
			Status: http.StatusInternalServerError, Error: err.Error()}, status.Errorf(codes.Internal, "Error during update: %s", err)
	}

	return &storepb.UpdateResponse{
		Status: http.StatusOK,
	}, nil
}
