package message

import (
	"github.com/PACZone/teleport/types/bypass"
	"github.com/PACZone/teleport/types/order"
)

type Message struct {
	To      bypass.Names
	From    bypass.Names
	Payload *order.Order
}
