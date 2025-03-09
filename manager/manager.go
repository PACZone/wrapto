package manager

import (
	"context"

	"github.com/PACZone/wrapto/config"
	"github.com/PACZone/wrapto/database"
	"github.com/PACZone/wrapto/http"
	logger "github.com/PACZone/wrapto/log"
	"github.com/PACZone/wrapto/sides/evm"
	"github.com/PACZone/wrapto/sides/pactus"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
)

type Manager struct {
	ctx      context.Context
	cancel   context.CancelFunc
	highway  chan message.Message
	bypasses map[bypass.Name]chan message.Message

	actors *actors
}

type actors struct {
	pactus  *pactus.Side
	polygon *evm.Side
	bsc     *evm.Side
	base    *evm.Side

	http *http.Server
}

// ! NEW EVM.
func NewManager(ctx context.Context, cancel context.CancelFunc,
	cfg *config.Config, db *database.Database,
) (*Manager, error) {
	highway := make(chan message.Message, 10)                  // TODO: what should we use as size?
	bypasses := make(map[bypass.Name]chan message.Message, 10) // TODO: what should we use as size?

	pactusCh := make(chan message.Message, 10)
	polygonCh := make(chan message.Message, 10)
	bscCh := make(chan message.Message, 10)
	baseCh := make(chan message.Message, 10)

	lastState, err := db.GetState() //nolint
	if err != nil {
		return nil, err
	}

	pactusSide, err := pactus.NewSide(ctx, highway, lastState.Pactus, pactusCh,
		&cfg.Pactus, db)
	if err != nil {
		return nil, err
	}

	polygonSide, err := evm.NewSide(ctx, highway, lastState.Polygon, polygonCh,
		cfg.Polygon, db, bypass.POLYGON)
	if err != nil {
		return nil, err
	}

	bscSide, err := evm.NewSide(ctx, highway, lastState.Bsc, bscCh,
		cfg.Bsc, db, bypass.BSC)
	if err != nil {
		return nil, err
	}

	baseSide, err := evm.NewSide(ctx, highway, lastState.Base, baseCh,
		cfg.Base, db, bypass.BASE)
	if err != nil {
		return nil, err
	}

	httpServer := http.NewHTTP(ctx, cfg.HTTPServer, db,
		highway, &cfg.Pactus, cfg.Polygon, cfg.Bsc, cfg.Base)

	actors := &actors{
		pactus:  pactusSide,
		polygon: polygonSide,
		bsc:     bscSide,
		base:    baseSide,

		http: httpServer,
	}

	bypasses[bypass.POLYGON] = polygonCh
	bypasses[bypass.PACTUS] = pactusCh
	bypasses[bypass.BSC] = bscCh
	bypasses[bypass.BASE] = baseCh

	return &Manager{
		ctx:      ctx,
		cancel:   cancel,
		highway:  highway,
		bypasses: bypasses,

		actors: actors,
	}, nil
}

func (m *Manager) Start() {
	logger.Info("manager actor spawned")

	go m.actors.pactus.Start()
	go m.actors.polygon.Start()
	go m.actors.bsc.Start()
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

func (m *Manager) routing(msg message.Message) error {
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

func (m *Manager) isRegistered(name bypass.Name) (chan message.Message, bool) {
	v, ok := m.bypasses[name]

	return v, ok
}
