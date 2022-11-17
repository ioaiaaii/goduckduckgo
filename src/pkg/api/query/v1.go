package query

import (
	"encoding/json"
	"fmt"
	"goduckduckgo/pkg/api"
	"goduckduckgo/pkg/duckduckgo"
	"net/http"

	route "goduckduckgo/pkg/server/http"
)

type queryRequest struct {
	Query string `json:"query"`
}

type QueryAPI struct {
	baseAPI *api.BaseAPI
}

// NewQueryAPI returns an initialized QueryAPI type.
func NewQueryAPI() *QueryAPI {
	return &QueryAPI{
		baseAPI: api.NewBaseAPI(),
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
		// RespondWithError(w, http.StatusBadRequest, err.Error())
		fmt.Println(err)
	}

	if d.More() {
		fmt.Println("Extraneous data in payload")
	}

	switch r.Method {
	case "GET":
		api.RespondWithError(w, http.StatusBadRequest, "Currently not supported")
		return
	case "POST":

		q, err := duckduckgo.NewDDGQuery("https://api.duckduckgo.com", request.Query)
		if err != nil {
			fmt.Println(err)
			api.RespondWithError(w, http.StatusBadRequest, err.Error())

		}
		err = q.Do()
		if err != nil {
			fmt.Println(err)
			api.RespondWithError(w, http.StatusBadRequest, err.Error())

		}
		api.RespondWithJSON(w, http.StatusAccepted, q.Payload())

	}

	// fmt.Println(q.Payload().Answer)
	// return &queryData{
	// 	ResultType: res.Value.Type(),
	// 	Result:     res.Value,
	// 	Stats:      qs,
	// }, res.Warnings, nil, qry.Close

	// switch r.Method {
	// case "GET":
	// 	// q, err := duckduckgo.NewDDGQuery("https://api.duckduckgo.com", request.Query)
	// 	// if err != nil {
	// 	// 	fmt.Println(err)
	// 	// 	return
	// 	// }
	// 	// q.Do()
	// 	// RespondWithJSON(w, http.StatusAccepted, q.Payload())
	// 	fmt.Println("lelos")
	// 	return nil, nil, nil, nil
	// case "POST":
	// 	q, err := duckduckgo.NewDDGQuery("https://api.duckduckgo.com", request.Query)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return nil, nil, nil, nil
	// 	}
	// 	q.Do()
	// 	// fmt.Println(q.Payload().Answer)
	// 	// return &queryData{
	// 	// 	ResultType: res.Value.Type(),
	// 	// 	Result:     res.Value,
	// 	// 	Stats:      qs,
	// 	// }, res.Warnings, nil, qry.Close

	// 	return q.Payload().Answer, nil, nil, nil

	// }
	//return q.Payload().Answer, nil, nil, nil

}
