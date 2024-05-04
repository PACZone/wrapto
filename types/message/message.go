package message

import (
	"fmt"

	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/order"
)

type Message struct {
	To      bypass.Name
	From    bypass.Name
	Payload *order.Order
}

func NewMessage(to, from bypass.Name, payload *order.Order) Message {
	return Message{
		To:      to,
		From:    from,
		Payload: payload,
	}
}

func (m Message) Validate(to bypass.Name) error {
	if m.To != to {
		return InvalidMessageError{
			Reason: fmt.Sprintf("invalid destination(to) value: %s", to),
		}
	}

	return nil
}
