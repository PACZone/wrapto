package pactus

import (
	"strings"

	"github.com/PACZone/wrapto/types/bypass"
)

type Dest struct {
	BypassName bypass.Name
	Addr       string
}

func ParseMemo(memo string) (*Dest, error) {
	splitMemo := strings.Split(memo, "@")

	if len(splitMemo) != 2 {
		return nil, InvalidMemoError{}
	}

	return &Dest{
		BypassName: bypass.Name(splitMemo[1]),
		Addr:       splitMemo[0],
	}, nil
}
