package core

import (
	"context"
	"sync"

	"github.com/PACZone/wrapto/config"
	"github.com/PACZone/wrapto/database"
	logger "github.com/PACZone/wrapto/log"
	"github.com/PACZone/wrapto/manager"
)

type Core struct {
	mgr *manager.Manager
}

func NewCore(ctx context.Context, cancel context.CancelFunc) (*Core, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		cancel()

		return nil, err
	}

	logger.InitGlobalLogger(&cfg.Logger)

	db, err := database.NewDB(cfg.Database.DSN)
	if err != nil {
		cancel()

		return nil, err
	}

	mgr, err := manager.NewManager(ctx, cancel, cfg, db)
	if err != nil {
		cancel()

		return nil, err
	}

	return &Core{
		mgr: mgr,
	}, nil
}

func (c *Core) Start() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		c.mgr.Start()
	}()

	wg.Wait()
}
