package params

import (
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/pactus-project/pactus/types/amount"
)

const (
	FeeFraction amount.Amount = 200 // Amount / 200 = 0.5%
	MinimumFee  amount.Amount = 1e9 // 1 PAC
	MaximumFee  amount.Amount = 5e9 // 5 PAC

	MainBypass = bypass.PACTUS
)
