package store

import (
	"goduckduckgo/pkg/store/storepb"

	"google.golang.org/grpc"
)

func RegisterStoreServer(storeSrv storepb.StoreServer) func(*grpc.Server) {
	return func(s *grpc.Server) {
		storepb.RegisterStoreServer(s, storeSrv)
	}
}
