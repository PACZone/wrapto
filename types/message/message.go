package message

import (
	"github.com/PACZone/teleport/types/bypass_name"
	"github.com/PACZone/teleport/types/order"
)

type Message struct {
	To      bypass_name.BypassName
	From    bypass_name.BypassName
	Payload *order.Order
}
