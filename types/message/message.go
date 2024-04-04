package message

import (
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/order"
)

type Msg struct {
	To      bypass.Name
	From    bypass.Name
	Payload *order.Order
}

func NewMsg(to bypass.Name, from bypass.Name, Payload *order.Order) *Msg {
	return &Msg{
		To:      to,
		From:    from,
		Payload: Payload,
	}
}
