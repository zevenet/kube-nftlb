package types

import (
	"io"
)

// Header is a custom struct that contains a key-value pair as strings
type Header struct {
	Key   string
	Value string
}

// Action defines the request type (GET, POST, ...)
type Action string

// Action types
const (
	// GET method
	GET = Action("GET")
	// POST method
	POST = Action("POST")
	// DELETE method
	DELETE = Action("DELETE")
)

// Payload is the body in a POST request, and nil otherwise
type Payload io.Reader

// Request is a custom struct that contains the data needed for a regular request
type Request struct {
	Header  *Header
	Action  Action
	URL     *URL
	Payload Payload
}

func (a Action) String() string {
	return string(a)
}
