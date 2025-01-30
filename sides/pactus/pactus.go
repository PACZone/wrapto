package pactus

import (
	"context"
	"sync"

	"github.com/PACZone/wrapto/database"
	logger "github.com/PACZone/wrapto/log"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
)

type Side struct {
	client   *Client
	listener *Listener
	bridge   Bridge
	highway  chan message.Message

	ctx context.Context
}

func NewSide(ctx context.Context,
	highway chan message.Message, startBlock uint32,
	bp chan message.Message, cfg *Config,
	db *database.Database,
) (*Side, error) {
	client, err := newClient(context.Background(), cfg.RPCNode) //nolint
	if err != nil {
		return nil, err
	}

	wallet, err := openWallet(cfg.WalletPath, cfg.WalletAddr, cfg.WalletPass)
	if err != nil {
		return nil, err
	}

	listener := newListener(ctx, client, bypass.PACTUS, highway, startBlock, cfg.LockAddr, db)
	bridge := newBridge(ctx, wallet, bp, bypass.PACTUS, db)

	return &Side{
		client:   client,
		listener: listener,
		highway:  highway,
		bridge:   bridge,

		ctx: ctx,
	}, nil
}

func (s *Side) Start() {
	logger.Info("pactus actor spawned")

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
			logger.Error("error starting listener", "actor", bypass.PACTUS, "err", err)
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
			logger.Error("error starting bridge", "actor", bypass.PACTUS, "err", err)
		}

		wg.Done()
	}()

	wg.Wait()

	logger.Info("pactus actor stopped")
}
