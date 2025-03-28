package http

import (
	"context"

	"github.com/PACZone/wrapto/database"
	"github.com/PACZone/wrapto/sides/evm"
	"github.com/PACZone/wrapto/sides/pactus"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo          *echo.Echo
	db            *database.Database
	ctx           context.Context
	cfg           Config
	polygonClient evm.Client
	bscClient     evm.Client
	baseClient    evm.Client
	pactusClient  pactus.Client
	highway       chan message.Message
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ! NEW EVM.
func NewHTTP(ctx context.Context, cfg Config, db *database.Database,
	highway chan message.Message, pacCfg *pactus.Config, polCfg, bscCfg, baseCfg evm.Config,
) *Server {
	app := echo.New()

	polClient, err := evm.NewPublicClient(polCfg.RPCNode, polCfg.ContractAddr, polCfg.ChainID)
	if err != nil {
		return nil
	}

	bscClient, err := evm.NewPublicClient(bscCfg.RPCNode, bscCfg.ContractAddr, bscCfg.ChainID)
	if err != nil {
		return nil
	}

	baseClient, err := evm.NewPublicClient(baseCfg.RPCNode, baseCfg.ContractAddr, baseCfg.ChainID)
	if err != nil {
		return nil
	}

	pacClient, err := pactus.NewClient(ctx, pacCfg.RPCNode)
	if err != nil {
		return nil
	}

	return &Server{
		echo:          app,
		db:            db,
		ctx:           ctx,
		cfg:           cfg,
		polygonClient: *polClient,
		pactusClient:  *pacClient,
		bscClient:     *bscClient,
		baseClient:    *baseClient,
		highway:       highway,
	}
}

func (s *Server) Start() {
	s.echo.Use(middleware.CORS())
	s.echo.GET("/rescan/:id", s.rescan)
	s.echo.GET("/state/latest", s.latestState)
	s.echo.GET("/state/stats", s.stats)
	s.echo.GET("/health", s.health)
	s.echo.GET("/transactions/recent", s.recentTxs)
	s.echo.GET("/search", s.searchTx)
	s.echo.GET("/announcement", s.announcement)

	err := s.echo.Start(s.cfg.Port)
	if err != nil {
		s.highway <- message.NewMessage(bypass.MANAGER, bypass.HTTP, nil)
	}

	<-s.ctx.Done()
}
