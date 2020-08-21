package types

import (
	"fmt"
)

// URL has different fields that match any regular URL
type URL struct {
	Protocol string
	IP       string
	Port     int
	Path     string
}

// SetPath
func (u *URL) SetPath(path string) {
	u.Path = path
}

func (u *URL) String() string {
	return fmt.Sprintf("%s://%s:%d/farms%s", u.Protocol, u.IP, u.Port, u.Path)
}
