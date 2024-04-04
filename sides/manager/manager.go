package manager

import (
	"context"

	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
)

type Mgr struct {
	Ctx      context.Context
	Highway  chan *message.Message
	Bypasses map[bypass.Name]chan *message.Message
}

func NewManager(ctx context.Context) *Mgr {
	return &Mgr{
		Ctx:      ctx,
		Highway:  make(chan *message.Message, 10),
		Bypasses: make(map[bypass.Name]chan *message.Message, 10),
	}
}

func (m *Mgr) Start() {
	for {
		select {
		case msg := <-m.Highway:
			err := m.routing(msg)
			if err != nil {
				continue
			}
		case <-m.Ctx.Done():
			return
		}
	}
}

func (m *Mgr) RegisterBypass(name bypass.Name, b chan *message.Message) error {
	_, ok := m.isRegistered(name)
	if !ok {
		m.Bypasses[name] = b

		return nil
	}

	return DupBypassError{BypassName: name}
}

func (m *Mgr) routing(msg *message.Message) error {
	b, ok := m.isRegistered(msg.To)
	if !ok {
		return BypassNotFoundError{BypassName: msg.To}
	}
	b <- msg

	return nil
}

func (m *Mgr) isRegistered(name bypass.Name) (chan *message.Message, bool) {
	v, ok := m.Bypasses[name]

	return v, ok
}
