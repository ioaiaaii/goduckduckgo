// New func initializes the Hub struct
package duckduckgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goduckduckgo/pkg/duckduckgo/typespb"
	"io"
	"net/http"
	"net/url"
	"time"
)

type DuckDuckGoQuery struct {
	query    string
	format   string
	queryURL string
	answer   *typespb.DuckDuckGoResponse
}

func NewDDGQuery(URL, query string) (*DuckDuckGoQuery, error) {
	var q DuckDuckGoQuery
	u, err := url.Parse(URL + "/")
	if err != nil {
		return nil, err
	}
	qQ := u.Query()
	qQ.Set("q", query)
	qQ.Set("format", "json")
	qQ.Set("t", "goduckduckgo")

	u.RawQuery = qQ.Encode()

	q.query = query
	q.format = "json"
	q.queryURL = u.String()
	return &q, nil

}

func (q *DuckDuckGoQuery) Payload() *typespb.QueryPayload {

	r := &typespb.QueryPayload{
		Answer: q.answer,
	}

	return r
}

// func debug(data []byte, err error) {
// 	if err == nil {
// 		fmt.Printf("%s\n\n", data)
// 	} else {
// 		log.Fatalf("%s\n\n", err)
// 	}
// }

func (q *DuckDuckGoQuery) Do() error {

	var response typespb.DuckDuckGoResponse

	//https://www.loginradius.com/blog/engineering/tune-the-go-http-client-for-high-performance/
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	req, err := http.NewRequest(http.MethodGet, q.queryURL, nil)

	// https://bugz.pythonanywhere.com/golang/Unexpected-EOF-golang-http-client-error
	req.Close = true

	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(req)

	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return err

	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return err
	}

	//https://blog.devgenius.io/to-unmarshal-or-to-decode-json-processing-in-go-explained-e92fab5b648f
	//decoder
	//use the json decoder on a read buffer initialized to body

	buf := bytes.NewBuffer(resBody)
	err = json.NewDecoder(buf).Decode(&response)

	//unmarshal
	// err = json.Unmarshal(resBody, &response)

	if err != nil {
		fmt.Printf("failed to decode response: %s", err.Error())
		//debug(httputil.DumpResponse(res, true))
		return err
	}

	q.answer = &response
	return nil
}
