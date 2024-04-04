package message

import (
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/order"
)

type Message struct {
	To      bypass.Name
	From    bypass.Name
	Payload *order.Order
}
