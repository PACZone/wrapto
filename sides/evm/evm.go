package evm

import (
	"context"
	"sync"

	"github.com/PACZone/wrapto/database"
	logger "github.com/PACZone/wrapto/log"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
)

type Side struct {
	client     *Client
	listener   *Listener
	bridge     Bridge
	db         *database.Database
	highway    chan message.Message
	bypass     chan message.Message
	bypassName bypass.Name

	ctx context.Context
}

func NewSide(ctx context.Context, highway chan message.Message, startOrder uint32,
	bp chan message.Message, env string, cfg Config, db *database.Database, bn bypass.Name,
) (*Side, error) {
	client, err := newClient(cfg.RPCNode, cfg.PrivateKey, cfg.ContractAddr, cfg.ChainID)
	if err != nil {
		return nil, err
	}

	listener := newListener(ctx, client, bn, highway, startOrder, db)
	bridge := newBridge(ctx, bp, bn, client, db)

	return &Side{
		client:     client,
		db:         db,
		listener:   listener,
		bridge:     bridge,
		highway:    highway,
		bypass:     bp,
		bypassName: bn,

		ctx: ctx,
	}, nil
}

func (s *Side) Start() {
	logger.Info("evm actor spawned", "chain", s.bypassName)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		err := s.listener.Start()
		if err != nil {
			s.highway <- message.Message{
				To:      bypass.MANAGER,
				From:    s.bypassName,
				Payload: nil,
			}
			logger.Error("error starting listener", "actor", s.bypassName, "err", err)
		}

		wg.Done()
	}()

	go func() {
		err := s.bridge.Start()
		if err != nil {
			s.highway <- message.Message{
				To:      bypass.MANAGER,
				From:    s.bypassName,
				Payload: nil,
			}
			logger.Error("error starting bridge", "actor", s.bypassName, "err", err)
		}

		wg.Done()
	}()

	wg.Wait()

	logger.Info("evm actor stopped", "chain", s.bypassName)
}
