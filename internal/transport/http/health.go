package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthServer struct {
	*echo.Echo
	host string
}

func NewHealthServer(host string) *HealthServer {
	return &HealthServer{
		Echo: echo.New(),
		host: host,
	}
}

func (hs *HealthServer) Run() error {
	hs.GET("/liveness", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	hs.GET("/readiness", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	hs.HideBanner = true

	err := hs.Start(hs.host)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
