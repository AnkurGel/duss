package server

import (
	"github.com/labstack/echo"
	"github.com/ankurgel/duss/internal/duss/models/auth"
	"strings"
	jwt "github.com/dgrijalva/jwt-go"
)

func (h *Handler) AuthenticateUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if skipTokenAuth(c) {
			return next(c)
		}
		user, err := GetAuthUser(h, c)
		if err != nil {
			c.Error(echo.ErrValidatorNotRegistered)
			return nil
		}
		if user.Email != c.Request().Header.Get("x-client-id") {
			c.Error(echo.ErrUnauthorized)
			return nil
		}
		return next(c)
	}
}

func GetAuthUser(h *Handler, c echo.Context) (*auth.User, error){
	user := c.Get("user").(*jwt.Token)
	authToken := user.Claims.(*auth.Token)
	authUser, err := h.Store.GetUserFromToken(authToken)
	if err != nil {
		return nil, err
	}
	return authUser, nil
}

func skipTokenAuth(c echo.Context) bool {
	if c.Path() != "/shorten" && !strings.Contains(c.Path(), "/admin"){
		return true
	}
	return false
}

func (h *Handler) VerifyAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := GetAuthUser(h, c)
		if err != nil {
			c.Error(echo.ErrValidatorNotRegistered)
			return nil
		}
		if user.Admin != 1 {
			c.Error(echo.ErrUnauthorized)
			return nil
		}
		return next(c)
	}
}
