package query

import (
	"context"
	"encoding/json"
	"fmt"
	"goduckduckgo/pkg/api"
	"goduckduckgo/pkg/config"
	"goduckduckgo/pkg/duckduckgo"
	"goduckduckgo/pkg/store/storepb"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	route "goduckduckgo/pkg/server/http"
)

type queryRequest struct {
	Query string `json:"query"`
}

type QueryAPI struct {
	baseAPI     *api.BaseAPI
	storeClient storepb.StoreClient
}

// NewQueryAPI returns an initialized QueryAPI type.
func NewQueryAPI(s storepb.StoreClient) *QueryAPI {
	return &QueryAPI{
		baseAPI:     api.NewBaseAPI(),
		storeClient: s,
	}
}

// Register the API's endpoints in the given router.
func (qapi *QueryAPI) Register(r *route.Router) {
	r.HandleRoute("/query", qapi.query)
	r.HandleRoute("/status/runtimeinfo", qapi.baseAPI.ServeRuntimeInfo)
	r.HandleRoute("/status/buildinfo", qapi.baseAPI.ServeBuildInfo)
}

func (qapi *QueryAPI) query(w http.ResponseWriter, r *http.Request) {
	var request *queryRequest

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(&request)

	// Handle payload status and values
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	if d.More() {
		api.RespondWithError(w, http.StatusBadRequest, "Extraneous data in payload")
	}

	switch r.Method {
	case "GET":
		api.RespondWithError(w, http.StatusBadRequest, "Currently not supported")
		return
	case "POST":

		q, err := duckduckgo.NewDDGQuery(config.DefaultDDGEndpoint, request.Query)
		if err != nil {
			fmt.Println(err)
			api.RespondWithError(w, http.StatusBadRequest, err.Error())

		}

		err = q.Do()
		if err != nil {
			fmt.Println(err)
			api.RespondWithError(w, http.StatusBadRequest, err.Error())

		}

		{
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			rr := storepb.CreateRequest{
				Query:  request.Query,
				Answer: q.Payload(),
			}

			fmt.Println(rr.Query)

			grpcReq, err := qapi.storeClient.Create(ctx, &rr)
			if err != nil {
				fmt.Println(err)
			}

			log.Info().Msgf("URL: %v", grpcReq.Status)
		}

		api.RespondWithJSON(w, http.StatusAccepted, q.Payload())

	}

}
