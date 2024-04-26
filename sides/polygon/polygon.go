package polygon

import (
	"context"
	"sync"

	"github.com/PACZone/wrapto/config"
	"github.com/PACZone/wrapto/database"
	logger "github.com/PACZone/wrapto/log"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
)

const (
	mainChainID   int64 = 137
	mumbaiChainID int64 = 80002
)

type Side struct {
	client   *Client
	listener *Listener
	bridge   Bridge
	db       *database.DB
	highway  chan message.Message
	bypass   chan message.Message

	ctx context.Context
}

func NewSide(ctx context.Context, highway chan message.Message, startOrder uint32,
	bp chan message.Message, env string, cfg config.PolygonConfig, db *database.DB,
) (*Side, error) {
	chainID := mainChainID
	if env == "dev" {
		chainID = mumbaiChainID
	}

	client, err := newClient(cfg.RPCNode, cfg.PrivateKey, cfg.ContractAddr, chainID)
	if err != nil {
		return nil, err
	}

	listener := newListener(ctx, client, bypass.POLYGON, highway, startOrder, db)
	bridge := newBridge(ctx, bp, bypass.POLYGON, client, db)

	return &Side{
		client:   client,
		db:       db,
		listener: listener,
		bridge:   bridge,
		highway:  highway,
		bypass:   bp,

		ctx: ctx,
	}, nil
}

func (s *Side) Start() {
	logger.Info("polygon actor spawned")

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		err := s.listener.Start()
		if err != nil {
			s.highway <- message.Message{
				To:      bypass.MANAGER,
				From:    s.bridge.bypassName,
				Payload: nil,
			}
			logger.Error("error starting listener", "actor", bypass.POLYGON, "err", err)
		}

		wg.Done()
	}()

	go func() {
		err := s.bridge.Start()
		if err != nil {
			s.highway <- message.Message{
				To:      bypass.MANAGER,
				From:    s.bridge.bypassName,
				Payload: nil,
			}
			logger.Error("error starting bridge", "actor", bypass.POLYGON, "err", err)
		}

		wg.Done()
	}()

	wg.Wait()

	logger.Info("polygon actor stopped")
}
