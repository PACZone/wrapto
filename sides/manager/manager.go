package manager

import (
	"context"

	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
)

type Mgr struct {
	ctx      context.Context
	highway  chan message.Message
	bypasses map[bypass.Name]chan message.Message
}

func NewManager(ctx context.Context) *Mgr {
	return &Mgr{
		ctx:      ctx,
		highway:  make(chan message.Message, 10),                 // TODO: what should we use as size?
		bypasses: make(map[bypass.Name]chan message.Message, 10), // TODO: what should we use as size?
	}
}

func (m *Mgr) Start() {
	for {
		select {
		case msg := <-m.highway:
			err := m.routing(msg)
			if err != nil {
				continue
			}
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *Mgr) RegisterBypass(name bypass.Name, b chan message.Message) error {
	_, ok := m.isRegistered(name)
	if !ok {
		m.bypasses[name] = b

		return nil
	}

	return DupBypassError{BypassName: name}
}

func (m *Mgr) routing(msg message.Message) error {
	b, ok := m.isRegistered(msg.To)
	if !ok {
		return BypassNotFoundError{BypassName: msg.To}
	}
	b <- msg

	return nil
}

func (m *Mgr) isRegistered(name bypass.Name) (chan message.Message, bool) {
	v, ok := m.bypasses[name]

	return v, ok
}
