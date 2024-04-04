package pactus

import "fmt"

type InvalidMemeError struct{}

func (e InvalidMemeError) Error() string {
	return "invalid memo"
}

type WalletNotExistError struct {
	path string
}

func (e WalletNotExistError) Error() string {
	return fmt.Sprintf("wallet not exist at: %s", e.path)
}
