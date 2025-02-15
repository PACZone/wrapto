package http

import (
	"fmt"
	"net/http"

	"github.com/PACZone/wrapto/types/bypass"
	"github.com/PACZone/wrapto/types/message"
	"github.com/PACZone/wrapto/types/order"
	"github.com/labstack/echo/v4"
)

// ! NEW EVM.
func (s *Server) rescan(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		res := Response{
			Status:  http.StatusBadRequest,
			Message: "Provide an ID please.",
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	ord, err := s.db.GetOrder(id)
	if err != nil {
		res := Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Can't find order with ID: %s.", id),
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	if ord.Status != order.FAILED {
		res := Response{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("Order status is: %s", ord.Status),
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	msg := message.Message{}

	switch ord.BridgeType {
	case order.PACTUS_POLYGON:
		msg = message.NewMessage(bypass.POLYGON, bypass.HTTP, ord)
	case order.POLYGON_PACTUS:
		msg = message.NewMessage(bypass.PACTUS, bypass.HTTP, ord)
	}

	s.highway <- msg

	res := Response{
		Status:  http.StatusOK,
		Message: "Ok",
	}

	return c.JSON(http.StatusOK, res)
}
