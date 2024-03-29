package aggregates

import "time"

// Event represents an event coming from browser/mobile, including both
// request and response data.
type Event struct {
	ID                string
	BrowserID         string
	ClientID          string
	Handled           interface{}
	ReplacesClientID  *string
	ResultingClientID string
	EventTime         time.Time
	Request           *Request
	Response          *Response

	// Processed Information
	ProgramIDs []string
}

// Request struct represents a browser/mobile request.
type Request struct {
	RequestTime    time.Time
	Body           *string
	BodyUsed       bool
	Cache          string
	Credentials    string
	Destination    string
	Headers        interface{}
	Integrity      string
	Method         string
	Mode           string
	Redirect       string
	Referrer       string
	ReferrerPolicy string
	URL            string
	Signal         interface{}
}

// Response struct represents a browser/mobile response.
type Response struct {
	ResponseTime time.Time
	Body         *string
	BodyUsed     bool
	Headers      interface{}
	Ok           bool
	Redirected   bool
	Status       uint16
	StatusText   string
	ResponseType string
	URL          string
}
