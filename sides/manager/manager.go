package manager

import (
	"context"

	"github.com/PACZone/teleport/types/bypass_name"
	"github.com/PACZone/teleport/types/message"
)

type Mgr struct {
	Ctx      context.Context
	Highway  chan *message.Message
	Bypasses map[bypass_name.BypassName]chan *message.Message
}

// new
func NewManager(ctx context.Context) *Mgr {
	return &Mgr{
		Ctx:      ctx,
		Highway:  make(chan *message.Message, 10),
		Bypasses: make(map[bypass_name.BypassName]chan *message.Message, 10),
	}
}

// start
func (m *Mgr) Start() {
	for {
		select {
		case msg := <-m.Highway:
			m.routing(msg)
		case <-m.Ctx.Done():
			return
		}
	}
}

// registerBypass
func (m *Mgr) RegisterBypass(name bypass_name.BypassName, b chan *message.Message) error {
	_, ok := m.isRegistered(name)
	if !ok {
		m.Bypasses[name] = b
		return nil
	}
	return DupBypassError{BypassName: name}
}

// routing
func (m *Mgr) routing(msg *message.Message) error {
	b, ok := m.isRegistered(msg.To)
	if !ok {
		return BypassNotFoundError{BypassName: msg.To}
	}
	b <- msg
	return nil
}

func (m *Mgr) isRegistered(name bypass_name.BypassName) (chan *message.Message, bool) {
	v, ok := m.Bypasses[name]
	return v, ok
}
