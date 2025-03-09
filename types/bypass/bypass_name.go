package bypass

type (
	Name string
)

// ! NEW EVM.
const (
	PACTUS  Name = "PACTUS"
	POLYGON Name = "POLYGON"
	BSC     Name = "BSC"
	BASE    Name = "BASE"
	MANAGER Name = "MANAGER"
	HTTP    Name = "HTTP"
)

func (n Name) ToStateName() string {
	switch n {
	case BSC:
		return "bsc"
	case PACTUS:
		return "pactus"
	case POLYGON:
		return "polygon"
	case BASE:
		return "base"
	case HTTP:
		return ""
	case MANAGER:
		return ""
	default:
		return ""
	}
}

// ! NEW EVM.
var ValidDestinations = []Name{POLYGON, BSC, BASE}
