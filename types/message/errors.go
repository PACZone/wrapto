package message

import "fmt"

type BasicCheckError struct {
	Reason string
}

func (e BasicCheckError) Error() string {
	return fmt.Sprintf("invalid message: %s", e.Reason)
}
