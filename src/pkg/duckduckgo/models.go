package duckduckgo

import (
	"net/url"
)

type Icon struct {
	//failed to decode response: json: cannot unmarshal string into Go struct field Icon.RelatedTopics.Icon.Height of type int
	Height int    `json:"Height,omitempty"`
	URL    string `json:"URL,omitempty"`
	Width  int    `json:"Width,omitempty"`
}

type Results struct {
	FirstURL string `json:"FirstURL,omitempty"`
	Icon     Icon   `json:"Icon,omitempty"`
	Result   string `json:"Result,omitempty"`
	Text     string `json:"Text,omitempty"`
}

type RelatedTopicsIcon struct {
	//failed to decode response: json: cannot unmarshal string into Go struct field Icon.RelatedTopics.Icon.Height of type int
	Height string `json:"Height,omitempty"`
	URL    string `json:"URL,omitempty"`
	Width  string `json:"Width,omitempty"`
}

type RelatedTopicsResults struct {
	FirstURL string            `json:"FirstURL,omitempty"`
	Icon     RelatedTopicsIcon `json:"Icon,omitempty"`
	Result   string            `json:"Result,omitempty"`
	Text     string            `json:"Text,omitempty"`
}

type DuckDuckGoResponse struct {
	Abstract         string
	AbstractSource   string
	AbstractText     string
	AbstractURL      string
	Answer           string `json:"Answer"`
	AnswerType       string `json:"AnswerType"`
	Definition       string `json:"Definition"`
	DefinitionSource string `json:"DefinitionSource"`
	DefinitionURL    string `json:"DefinitionURL"`
	Heading          string `json:"Heading"`
	Image            string `json:"Image"`
	ImageHeight      int    `json:"ImageHeight"`
	ImageIsLogo      int    `json:"ImageIsLogo"`
	ImageWidth       int    `json:"ImageWidth"`
	Redirect         string `json:"Redirect"`
	RelatedTopics    []RelatedTopicsResults
	Results          []Results
	Type             string
}

type DuckDuckGoQuery struct {
	query    string
	format   string
	queryURL string
	answer   *DuckDuckGoResponse
}

func NewDDGQuery(URL, query string) (*DuckDuckGoQuery, error) {
	var q DuckDuckGoQuery
	//Use a descriptive t parameter, i.e. append &t=nameofapp to your requests.
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

type QueryPayload struct {
	Answer DuckDuckGoResponse `json:"answer"`
}

func (q *DuckDuckGoQuery) Payload() *QueryPayload {
	r := &QueryPayload{
		Answer: *q.answer,
	}
	return r
}
