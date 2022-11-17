/*
Package api, provides a base API with an idiomatic way.
It contains methods and definitions that can be propagate to new APIs.
*/
package api

import (
	"encoding/json"
	"goduckduckgo/pkg/version"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type BaseAPI struct {
	runtimeInfo *version.GDDGRuntime
	buildInfo   *version.GDDGVersion
	Now         func() time.Time
}

// NewBaseAPI returns a new initialized BaseAPI type.
func NewBaseAPI() *BaseAPI {

	return &BaseAPI{
		runtimeInfo: version.RuntimeInfo,
		buildInfo:   version.BuildInfo,
		Now:         time.Now,
	}
}

func (api *BaseAPI) ServeRuntimeInfo(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusAccepted, api.runtimeInfo)
}

func (api *BaseAPI) ServeBuildInfo(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusAccepted, api.buildInfo)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	//Make sure to fetch fresh error message.
	//The no-store directive means browsers aren’t allowed to cache a response and must pull it from the server each time it’s requested.
	w.Header().Set("Cache-Control", "no-store")
	RespondWithJSON(w, code, map[string]string{"error": msg})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Warn().Err(err).Msgf("Responding with JSON Payload failed with error: %v", err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Warn().Err(err).Msgf("Responding with JSON Payload failed with error: %v", err.Error())

	}
}
