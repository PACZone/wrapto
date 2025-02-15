package bypass

type (
	Name string
)

// ! NEW EVM.
const (
	PACTUS  Name = "PACTUS"
	POLYGON Name = "POLYGON"
	MANAGER Name = "MANAGER"
	HTTP    Name = "HTTP"
)

var ValidDestinations = []Name{POLYGON}
