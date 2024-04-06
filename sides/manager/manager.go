package manager

import (
	"context"

	"github.com/PACZone/wrapto/config"
	"github.com/PACZone/wrapto/database"
	logger "github.com/PACZone/wrapto/log"
	"github.com/PACZone/wrapto/sides/pactus"
	"github.com/PACZone/wrapto/sides/polygon"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
)

type Mgr struct {
	ctx      context.Context
	cancel   context.CancelFunc
	highway  chan message.Message
	bypasses map[bypass.Name]chan message.Message

	sides *sides
}

type sides struct {
	pactus  *pactus.Side
	polygon *polygon.Side
}

func NewManager(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, db *database.DB) (*Mgr, error) {
	highway := make(chan message.Message, 10)                  // TODO: what should we use as size?
	bypasses := make(map[bypass.Name]chan message.Message, 10) // TODO: what should we use as size?

	pactusCh := make(chan message.Message, 10)
	polygonCh := make(chan message.Message, 10)

	pactusSide, err := pactus.NewSide(ctx, highway, 622508, pactusCh,
		cfg.Environment, cfg.Pactus, db) // TODO: retrieve the number from database.
	if err != nil {
		return nil, err
	}

	polygonSide, err := polygon.NewSide(ctx, highway, 6, polygonCh,
		cfg.Environment, cfg.Polygon, db) // TODO: retrieve the number from database.
	if err != nil {
		return nil, err
	}

	sides := &sides{
		pactus:  pactusSide,
		polygon: polygonSide,
	}

	bypasses[bypass.POLYGON] = polygonCh
	bypasses[bypass.PACTUS] = pactusCh

	return &Mgr{
		ctx:      ctx,
		cancel:   cancel,
		highway:  highway,
		bypasses: bypasses,

		sides: sides,
	}, nil
}

func (m *Mgr) Start() {
	logger.Info("manager actor spawned")

	go m.sides.pactus.Start()
	go m.sides.polygon.Start()

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

func (m *Mgr) routing(msg message.Message) error {
	if msg.To == bypass.MANAGER && msg.Payload == nil {
		m.cancel()

		return nil
	}

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
