package secret

import (
	"fmt"
	"regexp"
)

var (
	handlePattern = regexp.MustCompile(`^(?P<provider>[\w-]+):(?P<id>.+):(?P<key>.+)$`)
)

// Handle is the parsed representation of a secret handle from config
type Handle struct {
	Handle   string
	Provider string
	ID       string
	Key      string
}

// ParseHandle parses the string representation of a handle breaking it into its parts
func ParseHandle(h string) (handle *Handle, err error) {

	// extract parts from raw handle
	match := handlePattern.FindStringSubmatch(h)
	if match == nil {
		return nil, fmt.Errorf("unexpected handle format: %s", h)
	}

	// build secretHandle
	handle = &Handle{
		Handle:   h,
		Provider: match[1],
		ID:       match[2],
		Key:      match[3],
	}

	return handle, nil
}
