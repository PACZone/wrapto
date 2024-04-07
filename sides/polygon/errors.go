package polygon

import "fmt"

type ClientError struct {
	reason string
}

func (e ClientError) Error() string {
	return fmt.Sprintf("client error: %s", e.reason)
}
