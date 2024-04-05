package config

import "fmt"

type InvalidEnvironmentError struct {
	Environment string
}

func (e InvalidEnvironmentError) Error() string {
	return fmt.Sprintf("environment can be `dev` or `prod` not: %s", e.Environment)
}
