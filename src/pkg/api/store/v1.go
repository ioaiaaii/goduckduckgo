package store

import (
	"goduckduckgo/pkg/api"
)

type StoreAPI struct {
	baseAPI *api.BaseAPI
}

func NewStoreAPI() *StoreAPI {
	return &StoreAPI{
		baseAPI: api.NewBaseAPI(),
	}
}
