package pactus

type InvalidMemeError struct{}

func (e InvalidMemeError) Error() string {
	return "invalid memo"
}
