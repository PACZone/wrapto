package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Announcement struct {
	Title       string `json:"title"`
	Description string `json:"desc"`
	Show        bool   `json:"show"`
}

func (s *Server) announcement(c echo.Context) error {
	announc, err := s.db.GetAnnouncement()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  http.StatusInternalServerError,
			Message: "Internal error",
		})
	}

	res := Response{
		Status: http.StatusOK,
		Data: Announcement{
			Title:       announc.Title,
			Description: announc.Description,
			Show:        announc.Show,
		},
		Message: "Ok",
	}

	return c.JSON(http.StatusOK, res)
}
