package config

import "fmt"

type InvalidNetworkError struct {
	Network string
}

func (e InvalidNetworkError) Error() string {
	return fmt.Sprintf("network can be `main` or `test` not: %s", e.Network)
}
