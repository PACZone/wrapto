package core

import (
	"context"
	"sync"

	"github.com/PACZone/wrapto/config"
	"github.com/PACZone/wrapto/database"
	"github.com/PACZone/wrapto/sides/manager"
)

type Core struct {
	mgr *manager.Mgr
}

func NewCore() (*Core, error) {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.LoadConfig()
	if err != nil {
		cancel()

		return nil, err
	}

	db, err := database.NewDB(cfg.Database.Path)
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
