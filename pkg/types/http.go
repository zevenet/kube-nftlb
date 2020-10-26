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

// URL makes a formatted URL string with some parameters.
func URL(protocol string, host string, port int, path string) string {
	return fmt.Sprintf("%s://%s:%d/%s", protocol, host, port, path)
}
