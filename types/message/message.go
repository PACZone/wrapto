package message

import (
	"github.com/PACZone/teleport/types/bypass_name"
	"github.com/PACZone/teleport/types/order"
)

type Message struct {
	To      bypassname.BypassName
	From    bypassname.BypassName
	Payload *order.Order
}
