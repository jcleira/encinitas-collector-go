package handlers

import (
	"github.com/jcleira/encinitas-collector-go/internal/app/agent/aggregates"
)

// httpEventRequest represents the http version of an event coming from
// browser/mobile, including both request and response data.
type httpEventRequest struct {
	Event    httpEvent     `json:"event"`
	Request  *httpRequest  `json:"request,omitempty"`
	Response *httpResponse `json:"response,omitempty"`
}

// httpEvent represents the http version of an event coming from
// browser/mobile, including both request and response data.
type httpEvent struct {
	ID                string      `json:"id"`
	BrowserID         string      `json:"browser_id"`
	ClientID          string      `json:"client_id"`
	Handled           interface{} `json:"handled"`
	ReplacesClientID  *string     `json:"replaces_client_id,omitempty"`
	ResultingClientID string      `json:"resulting_client_id"`
	Duration          int64       `json:"duration"`
}

// httpRequest struct represents the http version of a browser/mobile request.
type httpRequest struct {
	Body           *string     `json:"body,omitempty"`
	BodyUsed       bool        `json:"body_used"`
	Cache          string      `json:"cache"`
	Credentials    string      `json:"credentials"`
	Destination    string      `json:"destination"`
	Headers        interface{} `json:"headers"`
	Integrity      string      `json:"integrity"`
	Method         string      `json:"method"`
	Mode           string      `json:"mode"`
	Redirect       string      `json:"redirect"`
	Referrer       string      `json:"referrer"`
	ReferrerPolicy string      `json:"referrer_policy"`
	URL            string      `json:"url"`
	Signal         interface{} `json:"signal"`
}

// httpResponse struct represents the http version of a browser/mobile
// response.
type httpResponse struct {
	Body         *string     `json:"body,omitempty"`
	BodyUsed     bool        `json:"body_used"`
	Headers      interface{} `json:"headers"`
	Ok           bool        `json:"ok"`
	Redirected   bool        `json:"redirected"`
	Status       uint16      `json:"status"`
	StatusText   string      `json:"status_text"`
	ResponseType string      `json:"response_type"`
	URL          string      `json:"url"`
}

func (her *httpEventRequest) ToAggregate() aggregates.Event {
	var request *aggregates.Request
	if her.Request != nil {
		request = her.Request.ToAggregate()
	}

	var response *aggregates.Response
	if her.Response != nil {
		response = her.Response.ToAggregate()
	}

	return aggregates.Event{
		ID:                her.Event.ID,
		BrowserID:         her.Event.BrowserID,
		ClientID:          her.Event.ClientID,
		Handled:           her.Event.Handled,
		ReplacesClientID:  her.Event.ReplacesClientID,
		ResultingClientID: her.Event.ResultingClientID,
		Duration:          her.Event.Duration,
		Request:           request,
		Response:          response,
	}
}

func (hr *httpRequest) ToAggregate() *aggregates.Request {
	return &aggregates.Request{
		Body:           hr.Body,
		BodyUsed:       hr.BodyUsed,
		Cache:          hr.Cache,
		Credentials:    hr.Credentials,
		Destination:    hr.Destination,
		Headers:        hr.Headers,
		Integrity:      hr.Integrity,
		Method:         hr.Method,
		Mode:           hr.Mode,
		Redirect:       hr.Redirect,
		Referrer:       hr.Referrer,
		ReferrerPolicy: hr.ReferrerPolicy,
		URL:            hr.URL,
		Signal:         hr.Signal,
	}
}

func (hr *httpResponse) ToAggregate() *aggregates.Response {
	return &aggregates.Response{
		Body:         hr.Body,
		BodyUsed:     hr.BodyUsed,
		Headers:      hr.Headers,
		Ok:           hr.Ok,
		Redirected:   hr.Redirected,
		Status:       hr.Status,
		StatusText:   hr.StatusText,
		ResponseType: hr.ResponseType,
		URL:          hr.URL,
	}
}
