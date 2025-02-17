package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type LatestStateResponse struct {
	Pactus  uint32 `json:"pactus"`
	Polygon uint32 `json:"polygon"`
}

type Stats struct {
	TotalSuccessfulBridges int     `json:"total_successful_bridges"`
	TotalWPACs             float64 `json:"total_wpacs"`
	TotalLockedPACs        float64 `json:"total_locked_pacs"`
}

func (s *Server) latestState(c echo.Context) error {
	state, err := s.db.GetState()
	if err != nil {
		res := Response{
			Status:  http.StatusInternalServerError,
			Message: "Error",
		}

		return c.JSON(http.StatusInternalServerError, res)
	}

	res := Response{
		Status:  http.StatusOK,
		Message: "Ok",
		Data: LatestStateResponse{
			Pactus:  state.Pactus,
			Polygon: state.Polygon,
		},
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) stats(c echo.Context) error {
	count, err := s.db.SuccessfulOrdersCount()
	if err != nil {
		res := Response{
			Status:  http.StatusInternalServerError,
			Message: "Error",
		}

		return c.JSON(http.StatusInternalServerError, res)
	}

	wpacCount, err := s.polygonClient.TotalSupply()
	if err != nil {
		res := Response{
			Status:  http.StatusInternalServerError,
			Message: "Error",
		}

		return c.JSON(http.StatusInternalServerError, res)
	}

	lockedPACsCount, err := s.pactusClient.GetTotalLocked()
	if err != nil {
		res := Response{
			Status:  http.StatusInternalServerError,
			Message: "Error",
		}

		return c.JSON(http.StatusInternalServerError, res)
	}

	res := Response{
		Status:  http.StatusOK,
		Message: "Ok",
		Data: Stats{
			TotalSuccessfulBridges: count,
			TotalWPACs:             wpacCount,
			TotalLockedPACs:        lockedPACsCount,
		},
	}

	return c.JSON(http.StatusOK, res)
}
