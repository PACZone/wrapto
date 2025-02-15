package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type SearchRequest struct {
	Q string `query:"q"`
}

func (s *Server) searchTx(c echo.Context) error {
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

	txs, err := s.db.SearchOrders(q.Q)
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
