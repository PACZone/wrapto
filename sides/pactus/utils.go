package pactus

import (
	"regexp"
	"slices"
	"strings"

	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/order"
)

type Dest struct {
	BypassName bypass.Name
	Addr       string
}

func ParseMemo(memo string) (Dest, error) {
	splitMemo := strings.Split(memo, "@")

	if len(splitMemo) != 2 {
		return Dest{}, InvalidMemoError{}
	}

	match, _ := regexp.MatchString("^(0x)?[0-9a-fA-F]{40}$", splitMemo[0])

	if !match {
		return Dest{}, InvalidMemoError{}
	}

	if !slices.Contains(bypass.ValidDestinations, bypass.Name(splitMemo[1])) {
		return Dest{}, InvalidMemoError{}
	}

	return Dest{
		BypassName: bypass.Name(splitMemo[1]),
		Addr:       splitMemo[0],
	}, nil
}

func (d *Dest) GetBridgeType() order.BridgeType {
	switch d.BypassName { //nolint
	case bypass.POLYGON:
		return order.PACTUS_POLYGON
	case bypass.BSC:
		return order.PACTUS_BSC
	}

	return ""
}
