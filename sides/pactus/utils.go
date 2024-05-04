package pactus

import (
	"regexp"
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

	match ,_ := regexp.MatchString("^(0x)?[0-9a-fA-F]{40}$",splitMemo[1])

	if(!match){
		return nil, InvalidMemoError{}
	}

	return &Dest{
		BypassName: bypass.Name(splitMemo[1]),
		Addr:       splitMemo[0],
	}, nil
}
