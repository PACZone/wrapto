package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/PACZone/wrapto/database"
	logger "github.com/PACZone/wrapto/log"
	"github.com/PACZone/wrapto/types/order"
	"github.com/labstack/echo/v4"
)

type HttpServer struct {
	echo *echo.Echo
	db   *database.DB
	ctx  context.Context
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
	Q  string `query:"q"`
}

type RecentTxsResponse struct {
	From   string        `json:"from"`
	To     string        `json:"to"`
	Amount float64 `json:"amount"`
	Fee    float64 `json:"fee"`
	Date   time.Time     `json:"date"`
	Status order.Status  `json:"status"`
	TxID   string        `json:"tx_id"`
}

func NewHttp(ctx context.Context, db *database.DB) *HttpServer {
	app := echo.New()

	return &HttpServer{
		echo: app,
		db:   db,
		ctx:  ctx,
	}
}

func (h *HttpServer) Start() {

	h.echo.GET("/state/latest", h.latestState)
	h.echo.GET("/health", h.health)
	h.echo.GET("/transactions/recent", h.recentTxs)
	h.echo.GET("/search", h.searchTx)

	h.echo.Start(":3000")
	logger.Info("http start on 3000")
	<-h.ctx.Done()
}

func (h *HttpServer) latestState(c echo.Context) (err error) {
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

func (h *HttpServer) health(c echo.Context) (err error) {
	res := Response{
		Status:  http.StatusOK,
		Message: "Ok",
		Data:    nil,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *HttpServer) recentTxs(c echo.Context) (err error) {

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
			From:   tx.Receiver,
			To:     tx.Sender,
			Fee:    tx.Fee.ToPAC(),
			Date:   tx.CreatedAt,
			Status: tx.Status,
			TxID:   tx.DestNetworkTxHash,
			Amount: tx.Amount.ToPAC(),
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

func (h *HttpServer) searchTx(c echo.Context) (err error) {

	var q SearchRequest

	err = c.Bind(&q); if err != nil {
		res := Response{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    nil,
		}
		return c.JSON(http.StatusBadRequest, res)
	}

	if strings.Trim(q.Q," ") == "" {
		res := Response{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    nil,
		}
		return c.JSON(http.StatusBadRequest, res)
	}

	txs,err:=h.db.SearchOrders(q.Q)
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
			From:   tx.Receiver,
			To:     tx.Sender,
			Fee:    tx.Fee.ToPAC(),
			Date:   tx.CreatedAt,
			Status: tx.Status,
			TxID:   tx.DestNetworkTxHash,
			Amount: tx.Amount.ToPAC(),
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
