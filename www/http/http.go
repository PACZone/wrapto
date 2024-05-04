package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/PACZone/wrapto/config"
	"github.com/PACZone/wrapto/database"
	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/PACZone/wrapto/types/order"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo    *echo.Echo
	db      *database.DB
	ctx     context.Context
	cfg     config.HTTPServerConfig
	highway chan message.Message
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type StateResponse struct {
	Pactus  uint32 `json:"pactus"`
	Polygon uint32 `json:"polygon"`
}

type SearchRequest struct {
	Q string `query:"q"`
}

type RecentTxsResponse struct {
	From       string       `json:"from"`
	To         string       `json:"to"`
	Amount     float64      `json:"amount"`
	Fee        float64      `json:"fee"`
	Date       time.Time    `json:"date"`
	Status     order.Status `json:"status"`
	TxID       string       `json:"tx_id"`
	BridgeType string       `json:"bridge_type"`
	Reason     string       `json:"reason"`
}

func NewHTTP(ctx context.Context, cfg config.HTTPServerConfig, db *database.DB, highway chan message.Message) *Server {
	app := echo.New()

	return &Server{
		echo:    app,
		db:      db,
		ctx:     ctx,
		cfg:     cfg,
		highway: highway,
	}
}

func (h *Server) Start() {
	h.echo.Use(middleware.CORS())
	h.echo.GET("/state/latest", h.latestState)
	h.echo.GET("/health", h.health)
	h.echo.GET("/transactions/recent", h.recentTxs)
	h.echo.GET("/search", h.searchTx)

	err := h.echo.Start(h.cfg.Port)
	if err != nil {
		h.highway <- message.NewMessage(bypass.MANAGER, bypass.HTTP, nil)
	}

	<-h.ctx.Done()
}

func (h *Server) latestState(c echo.Context) error {
	s, err := h.db.GetState()
	if err != nil {
		res := Response{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    nil,
		}

		return c.JSON(http.StatusInternalServerError, res)
	}
	res := Response{
		Status:  http.StatusOK,
		Message: "Ok",
		Data: StateResponse{
			Pactus:  s.Pactus,
			Polygon: s.Polygon,
		},
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Server) health(c echo.Context) error {
	res := Response{
		Status:  http.StatusOK,
		Message: "Ok",
		Data:    nil,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Server) recentTxs(c echo.Context) error {
	txs, err := h.db.GetLatestOrders(10)
	if err != nil {
		res := Response{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    nil,
		}

		return c.JSON(http.StatusInternalServerError, res)
	}

	dto := make([]RecentTxsResponse, 0)

	for _, tx := range txs {
		a := RecentTxsResponse{
			From:       tx.Receiver,
			To:         tx.Sender,
			Fee:        tx.Fee.ToPAC(),
			Date:       tx.CreatedAt,
			Status:     tx.Status,
			TxID:       tx.DestNetworkTxHash,
			Amount:     tx.Amount.ToPAC(),
			BridgeType: string(tx.BridgeType),
			Reason:     tx.Reason,
		}
		dto = append(dto, a)
	}

	res := Response{
		Status:  http.StatusOK,
		Message: "Ok",
		Data:    dto,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Server) searchTx(c echo.Context) error {
	var q SearchRequest

	err := c.Bind(&q)
	if err != nil {
		res := Response{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    nil,
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	if strings.Trim(q.Q, " ") == "" {
		res := Response{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    nil,
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	txs, err := h.db.SearchOrders(q.Q)
	if err != nil {
		res := Response{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    nil,
		}

		return c.JSON(http.StatusInternalServerError, res)
	}

	dto := make([]RecentTxsResponse, 0)

	for _, tx := range txs {
		a := RecentTxsResponse{
			From:       tx.Receiver,
			To:         tx.Sender,
			Fee:        tx.Fee.ToPAC(),
			Date:       tx.CreatedAt,
			Status:     tx.Status,
			TxID:       tx.DestNetworkTxHash,
			Amount:     tx.Amount.ToPAC(),
			BridgeType: string(tx.BridgeType),
			Reason:     tx.Reason,
		}
		dto = append(dto, a)
	}

	res := Response{
		Status:  http.StatusOK,
		Message: "Ok",
		Data:    dto,
	}

	return c.JSON(http.StatusOK, res)
}
