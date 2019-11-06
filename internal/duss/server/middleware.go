package server

import (
	"github.com/ankurgel/duss/internal/duss/models/auth"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"strings"
)

// AuthenticateUser is a middleware to authenticate user
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

// GetAuthUser returns the Autherized user from token
func GetAuthUser(h *Handler, c echo.Context) (*auth.User, error){
	user := c.Get("user").(*jwt.Token)
	authToken := user.Claims.(*auth.Token)
	authUser, err := h.Store.GetUserFromToken(authToken)
	if err != nil {
		return nil, err
	}
	return authUser, nil
}

// skipTokenAuth is use to add paths to skip Token auth check
func skipTokenAuth(c echo.Context) bool {
	if c.Path() != "/shorten" && !strings.Contains(c.Path(), "/admin"){
		return true
	}
	return false
}

// VerifyAdmin is a middleware verify the admin role of the authenticate user
func (h *Handler) VerifyAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := GetAuthUser(h, c)
		if err != nil {
			c.Error(echo.ErrValidatorNotRegistered)
			return nil
		}
		if !user.IsAdmin() {
			c.Error(echo.ErrUnauthorized)
			return nil
		}
		return next(c)
	}
}
