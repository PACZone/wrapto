package manager

import (
	"fmt"

	"github.com/PACZone/teleport/types/bypass"
)

type DupBypassError struct {
	BypassName bypass.Names
}

func (e DupBypassError) Error() string {
	return fmt.Sprintf("%s is a duplicated bypass, you can add a bypass only once", e.BypassName)
}

type BypassNotFoundError struct {
	BypassName bypass.Names
}

func (e BypassNotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.BypassName)
}
