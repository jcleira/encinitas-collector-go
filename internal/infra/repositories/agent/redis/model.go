package redis

import (
	"github.com/jcleira/encinitas-collector-go/internal/app/agent/aggregates"
)

// redisEvent represents the redis version of an event coming from
// browser/mobile, including both request and response data.
type redisEvent struct {
	ID                string         `json:"id"`
	ClientID          string         `json:"client_id"`
	BrowserID         string         `json:"browser_id"`
	Handled           interface{}    `json:"handled"`
	ReplacesClientID  *string        `json:"replaces_client_id,omitempty"`
	ResultingClientID string         `json:"resulting_client_id"`
	Duration          int64          `json:"duration"`
	Request           *redisRequest  `json:"request,omitempty"`
	Response          *redisResponse `json:"response,omitempty"`
}

// redisRequest struct represents the redis version of a browser/mobile request.
type redisRequest struct {
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

// redisResponse struct represents the redis version of a browser/mobile
// response.
type redisResponse struct {
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

func (r *redisEvent) toAggregate() aggregates.Event {
	return aggregates.Event{
		ID:                r.ID,
		ClientID:          r.ClientID,
		BrowserID:         r.BrowserID,
		Handled:           r.Handled,
		ReplacesClientID:  r.ReplacesClientID,
		ResultingClientID: r.ResultingClientID,
		Request:           r.Request.toAggregate(),
		Response:          r.Response.toAggregate(),
	}
}

func (r *redisRequest) toAggregate() *aggregates.Request {
	return &aggregates.Request{
		Body:           r.Body,
		BodyUsed:       r.BodyUsed,
		Cache:          r.Cache,
		Credentials:    r.Credentials,
		Destination:    r.Destination,
		Headers:        r.Headers,
		Integrity:      r.Integrity,
		Method:         r.Method,
		Mode:           r.Mode,
		Redirect:       r.Redirect,
		Referrer:       r.Referrer,
		ReferrerPolicy: r.ReferrerPolicy,
		URL:            r.URL,
		Signal:         r.Signal,
	}
}

func (r *redisResponse) toAggregate() *aggregates.Response {
	return &aggregates.Response{
		Body:         r.Body,
		BodyUsed:     r.BodyUsed,
		Headers:      r.Headers,
		Ok:           r.Ok,
		Redirected:   r.Redirected,
		Status:       r.Status,
		StatusText:   r.StatusText,
		ResponseType: r.ResponseType,
		URL:          r.URL,
	}
}

func redisEventFromAggregate(event aggregates.Event) redisEvent {
	var redisRequest *redisRequest
	if event.Request != nil {
		redisRequest = redisRequestFromAggregate(*event.Request)
	}

	var redisResponse *redisResponse
	if event.Response != nil {
		redisResponse = redisResponseFromAggregate(*event.Response)
	}

	return redisEvent{
		ID:                event.ID,
		BrowserID:         event.BrowserID,
		ClientID:          event.ClientID,
		Handled:           event.Handled,
		ReplacesClientID:  event.ReplacesClientID,
		ResultingClientID: event.ResultingClientID,
		Request:           redisRequest,
		Response:          redisResponse,
	}
}

func redisRequestFromAggregate(request aggregates.Request) *redisRequest {
	return &redisRequest{
		Body:           request.Body,
		BodyUsed:       request.BodyUsed,
		Cache:          request.Cache,
		Credentials:    request.Credentials,
		Destination:    request.Destination,
		Headers:        request.Headers,
		Integrity:      request.Integrity,
		Method:         request.Method,
		Mode:           request.Mode,
		Redirect:       request.Redirect,
		Referrer:       request.Referrer,
		ReferrerPolicy: request.ReferrerPolicy,
		URL:            request.URL,
		Signal:         request.Signal,
	}
}

func redisResponseFromAggregate(response aggregates.Response) *redisResponse {
	return &redisResponse{
		Body:         response.Body,
		BodyUsed:     response.BodyUsed,
		Headers:      response.Headers,
		Ok:           response.Ok,
		Redirected:   response.Redirected,
		Status:       response.Status,
		StatusText:   response.StatusText,
		ResponseType: response.ResponseType,
		URL:          response.URL,
	}
}
