package message

import (
	"github.com/PACZone/teleport/types/bypass"
	"github.com/PACZone/teleport/types/order"
)

type Message struct {
	To      bypass.Name
	From    bypass.Name
	Payload *order.Order
}
