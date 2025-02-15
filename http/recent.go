package http

import (
	"net/http"
	"time"

	"github.com/PACZone/wrapto/types/order"
	"github.com/labstack/echo/v4"
)

type RecentTxsResponse struct {
	ID         string       `json:"id"`
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

func (s *Server) recentTxs(c echo.Context) error {
	txs, err := s.db.GetLatestOrders(10)
	if err != nil {
		res := Response{
			Status:  http.StatusInternalServerError,
			Message: "Error",
		}

		return c.JSON(http.StatusInternalServerError, res)
	}

	dto := make([]RecentTxsResponse, 0)

	for _, tx := range txs {
		a := RecentTxsResponse{
			ID:         tx.ID,
			From:       tx.Receiver,
			To:         tx.Sender,
			Fee:        tx.Fee().ToPAC(),
			Date:       tx.CreatedAt,
			Status:     tx.Status,
			TxID:       tx.DestNetworkTxHash,
			Amount:     tx.AmountAfterFee().ToPAC(),
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
