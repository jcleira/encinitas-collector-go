package handlers

import (
	"time"

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
	BrowserID         string      `json:"browserId"`
	ClientID          string      `json:"clientId"`
	Handled           interface{} `json:"handled"`
	ReplacesClientID  *string     `json:"replacesClientId,omitempty"`
	ResultingClientID string      `json:"resultingClientId"`
	EventTime         int64       `json:"eventTime"`
}

// httpRequest struct represents the http version of a browser/mobile request.
type httpRequest struct {
	RequestTime    int64       `json:"requestTime"`
	Body           *string     `json:"body,omitempty"`
	BodyUsed       bool        `json:"bodyUsed"`
	Cache          string      `json:"cache"`
	Credentials    string      `json:"credentials"`
	Destination    string      `json:"destination"`
	Headers        interface{} `json:"headers"`
	Integrity      string      `json:"integrity"`
	Method         string      `json:"method"`
	Mode           string      `json:"mode"`
	Redirect       string      `json:"redirect"`
	Referrer       string      `json:"referrer"`
	ReferrerPolicy string      `json:"referrerPolicy"`
	URL            string      `json:"url"`
	Signal         interface{} `json:"signal"`
}

// httpResponse struct represents the http version of a browser/mobile
// response.
type httpResponse struct {
	ResponseTime int64       `json:"responseTime"`
	Body         *string     `json:"body,omitempty"`
	BodyUsed     bool        `json:"bodyUsed"`
	Headers      interface{} `json:"headers"`
	Ok           bool        `json:"ok"`
	Redirected   bool        `json:"redirected"`
	Status       uint16      `json:"status"`
	StatusText   string      `json:"statusText"`
	ResponseType string      `json:"responseType"`
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
		EventTime:         time.Unix(0, her.Event.EventTime*int64(time.Millisecond)),
		Request:           request,
		Response:          response,
	}
}

func (hr *httpRequest) ToAggregate() *aggregates.Request {

	return &aggregates.Request{
		RequestTime:    time.Unix(0, hr.RequestTime*int64(time.Millisecond)),
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
		ResponseTime: time.Unix(0, hr.ResponseTime*int64(time.Millisecond)),
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
