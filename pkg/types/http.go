package types

import (
	"fmt"
	"io"
)

// RequestData contains the data needed for a regular request.
type RequestData struct {
	Method string
	Path   string
	Body   io.Reader
}

// URL has different fields that match any regular URL.
type URL struct {
	Protocol string
	IP       string
	Port     int
	Path     string
}

func (u *URL) String() string {
	return fmt.Sprintf("%s://%s:%d/farms%s", u.Protocol, u.IP, u.Port, u.Path)
}
