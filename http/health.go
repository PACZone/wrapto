package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (*Server) health(c echo.Context) error {
	res := Response{
		Status:  http.StatusOK,
		Message: "Ok",
	}

	return c.JSON(http.StatusOK, res)
}
