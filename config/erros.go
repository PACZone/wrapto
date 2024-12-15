package config

import "fmt"

type InvalidEnvironmentError struct {
	Environment string
}

func (e InvalidEnvironmentError) Error() string {
	return fmt.Sprintf("environment can be `dev` or `prod` not: %s", e.Environment)
}

// Error represents an error in loading or validating config.
type Error struct {
	reason string
}

func (e Error) Error() string {
	return fmt.Sprintf("config error: %s\n", e.reason)
}
