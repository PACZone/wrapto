package order

import "fmt"

type BasicCheckError struct {
	Reason string
}

func (e BasicCheckError) Error() string {
	return fmt.Sprintf("invalid bridge order: %s", e.Reason)
}
