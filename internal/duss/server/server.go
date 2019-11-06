// Package server contains the abstraction
// for the server and data-store instantiation and interaction
package server

import (
	"context"
	"fmt"
	"github.com/ankurgel/duss/internal/duss/algo"
	"github.com/ankurgel/duss/internal/duss/models/auth"
	"github.com/ankurgel/duss/internal/duss/models/url"
	"github.com/ankurgel/duss/internal/duss/store"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

// Handler is a manager for a storage and a router.
// It abstracts all the routes and their interaction with data
type Handler struct {
	Store  *store.Store
	Router *echo.Echo
}

// Handlers

// SetHandlers defines the routes and verb allowed in application
func (h *Handler) SetHandlers() {
	h.Router.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &auth.Token{},
		SigningKey: []byte(viper.GetString("JwtSecret")),
		Skipper: skipTokenAuth,
	}))

	h.Router.Use(h.AuthenticateUser)
	h.Router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "DUSS ka dum!")
	})

	h.Router.GET("/:shortUrl", func(c echo.Context) error {
		return getLongURL(h, c)
	})

	h.Router.POST("/shorten", func(c echo.Context) error {
		return cutShort(h, c)
	})

	adminGroup := h.Router.Group("/admin")
	adminGroup.Use(h.VerifyAdmin)
	
	adminGroup.POST("/admin/user/new", func(c echo.Context) error {
		return createNewUser(h, c)
	})
}

// cutShort responds with shortened URL or gives 422 if can't be processed
func cutShort(h *Handler, c echo.Context) error {
	var u string
	var e error
	custom := c.FormValue("custom")
	// TODO: move this in CreateByLongURL
	if u, e = algo.NormalizeURL(c.FormValue("url")); e != nil {
		errorMessage := fmt.Sprintf("Error in URL for %s: %s", c.FormValue("url"), e)
		log.Error(errorMessage)
		return c.String(http.StatusUnprocessableEntity, errorMessage)
	}
	result, e := h.Store.CreateByLongURL(u, custom)
	if e != nil {
		errorMessage := fmt.Sprintf("Error in shortening for %s: %s", c.FormValue("url"), e)
		log.Error(errorMessage)
		return c.String(http.StatusUnprocessableEntity, errorMessage)
	}
	return c.String(http.StatusCreated, result.ShortURL())
}

// getLongURL responds with original long URL redirect for a given short slug
func getLongURL(h *Handler, c echo.Context) error {
	shortURL := c.Param("shortUrl")
	var u *url.URL
	var e error
	if u, e = h.Store.FindByShortURL(shortURL); e != nil {
		log.Error(fmt.Sprintf("Error in getSlug for %s: %s", shortURL, e))
		return c.String(http.StatusNotFound, "Invalid Link")
	}
	return c.Redirect(http.StatusMovedPermanently, u.Original)
}


func createNewUser(h *Handler, c echo.Context) error {
	newUser := &auth.User{}
	newUser.Email = c.FormValue("email_id")
	newUser.Password = c.FormValue("password")
	newUser.Name = c.FormValue("name")
	var token string
	var e error
	if token, e = h.Store.CreateUser(newUser); e != nil {
		log.Error(fmt.Sprintf("Error in Create user for %s: %s", newUser.Email, e))
		return c.String(http.StatusNotFound, "Invalid Link")
	}
	newUser.Password = ""
	newUser.Token = token
	return c.JSON(http.StatusCreated, newUser)
}


// Start-Stop

// Listen setups the middlewares and listens on configured port
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

// Close gracefully shuts down the server on interrupt
func (h *Handler) Close() error {
	log.Info("Shutting down the server... ")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.Router.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// InitServer instantiates the server with the given data store
func InitServer(store *store.Store) *Handler {
	h := &Handler{Store: store}
	h.Router = echo.New()

	return h
}
