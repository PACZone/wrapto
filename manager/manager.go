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
	"github.com/PACZone/wrapto/www/http"
)

type Mgr struct {
	ctx      context.Context
	cancel   context.CancelFunc
	highway  chan message.Message
	bypasses map[bypass.Name]chan message.Message

	actors *actors
}

type actors struct {
	pactus  *pactus.Side
	polygon *polygon.Side

	http *http.Server
}

func NewManager(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, db *database.DB) (*Mgr, error) {
	highway := make(chan message.Message, 10)                  // TODO: what should we use as size?
	bypasses := make(map[bypass.Name]chan message.Message, 10) // TODO: what should we use as size?

	pactusCh := make(chan message.Message, 10)
	polygonCh := make(chan message.Message, 10)

	lastState, err := db.GetState()
	if err != nil {
		return nil, err
	}

	pactusSide, err := pactus.NewSide(ctx, highway, lastState.Pactus, pactusCh,
		cfg.Environment, cfg.Pactus, db)
	if err != nil {
		return nil, err
	}

	polygonSide, err := polygon.NewSide(ctx, highway, lastState.Polygon, polygonCh,
		cfg.Environment, cfg.Polygon, db)
	if err != nil {
		return nil, err
	}

	httpSide := http.NewHTTP(ctx, cfg.HTTPServer, db, highway)

	actors := &actors{
		pactus:  pactusSide,
		polygon: polygonSide,

		http: httpSide,
	}

	bypasses[bypass.POLYGON] = polygonCh
	bypasses[bypass.PACTUS] = pactusCh

	return &Mgr{
		ctx:      ctx,
		cancel:   cancel,
		highway:  highway,
		bypasses: bypasses,

		actors: actors,
	}, nil
}

func (m *Mgr) Start() {
	logger.Info("manager actor spawned")

	go m.actors.pactus.Start()
	go m.actors.polygon.Start()
	go m.actors.http.Start()

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
