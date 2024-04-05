package manager

import (
	"fmt"

	"github.com/PACZone/wrapto/types/bypass"
)

type BypassNotFoundError struct {
	BypassName bypass.Name
}

func (e BypassNotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.BypassName)
}
