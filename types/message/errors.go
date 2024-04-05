package message

import "fmt"

type InvalidMessageError struct {
	Reason string
}

func (e InvalidMessageError) Error() string {
	return fmt.Sprintf("invalid message: %s", e.Reason)
}
