package server

import (
	"context"
	"fmt"
	"github.com/ankurgel/duss/internal/duss/algo"
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

	h.Router.GET("/:shortUrl", func(c echo.Context) error {
		return getLongUrl(h, c)
	})

	h.Router.POST("/shorten", func(c echo.Context) error {
		return cutShort(h, c)
	})
}

func cutShort(h *Handler, c echo.Context) error {
	var u string
	var e error
	custom := c.FormValue("custom")
	// TODO: move this in CreateByLongUrl
	if u, e = algo.NormalizeUrl(c.FormValue("url")); e != nil {
		errorMessage := fmt.Sprintf("Error in URL for %s: %s", c.FormValue("url"), e)
		log.Error(errorMessage)
		return c.String(http.StatusUnprocessableEntity, errorMessage)
	}
	if result, e := h.Store.CreateByLongUrl(u, custom); e != nil {
		errorMessage := fmt.Sprintf("Error in shortening for %s: %s", c.FormValue("url"), e)
		log.Error(errorMessage)
		return c.String(http.StatusUnprocessableEntity, errorMessage)
	} else {
		return c.String(http.StatusCreated, result.ShortUrl())
	}
	return c.String(http.StatusOK, u)
}

func getLongUrl(h *Handler, c echo.Context) error {
	shortUrl := c.Param("shortUrl")
	var u *url.Url
	var e error
	if u, e = h.Store.FindByShortUrl(shortUrl); e != nil {
		log.Error(fmt.Sprintf("Error in getSlug for %s: %s", shortUrl, e))
		return c.String(http.StatusNotFound, "Invalid Link")
	}
	return c.Redirect(http.StatusMovedPermanently, u.Original)
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
