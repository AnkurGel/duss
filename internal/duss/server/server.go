package server

import (
	"context"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	log "github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	router *echo.Echo
}

// Handlers

func (h *Handler) SetHandlers() {
	h.router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "DUSS ka dum!")
	})

	h.router.GET("/:slug", getSlug)
	h.router.POST("/shorten", cutShort)
}

func cutShort(c echo.Context) error {
	url := c.FormValue("url")
	return c.String(http.StatusOK, url)
}

func getSlug(c echo.Context) error {
	slug := c.Param("slug")
	return c.String(http.StatusOK, slug)
}


// Start-Stop

func (h *Handler) Listen(port string) error {
	h.router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "HTTP ${method} ${uri} Response=${status} ${latency_human}\n",
	}))
	//h.router.Use(middleware.Logger())

	if err := h.router.Start(port); err != nil {
		return err
	}
	return nil
}

func (h *Handler) Close() error {
	log.Info("Shutting down the server... ")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.router.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func InitServer() *Handler {
	h := &Handler{}
	h.router = echo.New()

	return h
}
