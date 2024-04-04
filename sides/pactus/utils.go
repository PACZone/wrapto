package pactus

import (
	"strings"

	"github.com/PACZone/wrapto/types/bypass"
)

type dest struct {
	name bypass.Name
	addr string
}

func parseMemo(memo string) (*dest, error) {
	splitMemo := strings.Split(memo, "@")

	if len(splitMemo) != 2 {
		return nil, InvalidMemeError{}
	}

	return &dest{
		name: bypass.Name(splitMemo[1]),
		addr: splitMemo[0],
	}, nil
}
