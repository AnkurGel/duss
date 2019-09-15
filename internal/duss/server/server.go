package server

import (
	"context"
	"fmt"
	"github.com/ankurgel/duss/internal/duss/models/url"
	"github.com/ankurgel/duss/internal/duss/store"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Handler struct {
	Store  *store.Store
	Router *echo.Echo
}

// Handlers

func (h *Handler) SetHandlers() {
	h.Router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "DUSS ka dum!")
	})
	h.Router.GET("/:slug", func(c echo.Context) error {
		return getSlug(h, c)
	})
	//h.Router.GET("/:slug", getSlug)
	h.Router.POST("/shorten", cutShort)
}

func cutShort(c echo.Context) error {
	url := c.FormValue("url")
	//custom := c.FormValue("custom")
	return c.String(http.StatusOK, url)
}

func getSlug(h *Handler, c echo.Context) error {
	slug := c.Param("slug")
	var u *url.Url
	var e error
	if u, e = h.Store.FindBySlug(slug); e != nil {
		log.Print(fmt.Sprintf("Error in getSlug for %s: %s", slug, e))
		return c.String(http.StatusNotFound, slug)
	}
	return c.String(http.StatusOK, u.Original)
}


// Start-Stop

func (h *Handler) Listen(port string) error {
	h.Router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "HTTP ${method} ${uri} Response=${status} ${latency_human}\n",
	}))
	//h.Router.Use(middleware.Logger())

	if err := h.Router.Start(port); err != nil {
		return err
	}
	return nil
}

func (h *Handler) Close() error {
	log.Info("Shutting down the server... ")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.Router.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func InitServer(store *store.Store) *Handler {
	h := &Handler{Store: store}
	h.Router = echo.New()

	return h
}
