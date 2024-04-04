package types

import "github.com/PACZone/wrapto/types/bypass"

const (
	FeeFraction float64 = 0.005 // 0.5%
	MinimumFee  float64 = 1e9   // 1 PAC
	MaximumFee  float64 = 5e9   // 5 PAC
  
	MainBypass  bypass.Name = bypass.PACTUS
)
