package handler

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/config"
	"itisadb/internal/cli/cookies"
	"itisadb/internal/cli/schema"
	"itisadb/internal/cli/usecase"
	"net/http"
)

type Handler struct {
	logic    *usecase.UseCase
	security config.SecurityConfig
	*zap.Logger
}

func New(logic *usecase.UseCase, loggerInstance *zap.Logger, security config.SecurityConfig) *Handler {
	return &Handler{logic: logic, Logger: loggerInstance, security: security}
}

func (h *Handler) MainPage(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if err != nil || cookie == nil {
		return c.Redirect(http.StatusMovedPermanently, "/auth")
	}

	return c.Render(http.StatusOK, "index-2.html", nil)
}

func (h *Handler) GetAuthPage(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if err == nil && cookie != nil {
		return c.Redirect(http.StatusMovedPermanently, "/")
	}

	if !h.security.On || !h.security.MandatoryAuthorization {
		cookie = cookies.SetCookie("itisadb")
		c.SetCookie(cookie)
	}

	return c.Render(http.StatusOK, "auth.html", nil)
}

func (h *Handler) Action(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	cookie, err := c.Cookie("session")
	if err != nil || cookie == nil {
		return c.Redirect(http.StatusMovedPermanently, "/auth")
	}

	action := c.Request().URL.Query().Get("action")
	res, err := h.logic.ProcessQuery(c.Request().Context(), cookie.Value, action)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code().String() == codes.NotFound.String() {
			err = errors.New("not found")
		} else if st.Message() == "unknown server" {
			err = errors.New("unknown server")
		} else if ok && st.Code().String() == codes.Unavailable.String() {
			err = errors.New("memory balancer is offline")
		}

		var t = schema.Response{Text: err.Error()}
		bytes, err := json.Marshal(t)
		if err != nil {
			h.Warn(err.Error())
			return err
		}
		c.Response().Write(bytes)
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	var t = schema.Response{Text: res}

	bytes, err := json.Marshal(t)
	if err != nil {
		t = schema.Response{Text: err.Error()}
		bytes, err = json.Marshal(t)
		if err != nil {
			h.Warn(err.Error())
		}
	}
	c.Response().Write(bytes)
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

func (h *Handler) History(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	cookie, err := c.Cookie("session")
	if err != nil || cookie == nil {
		return c.Redirect(http.StatusMovedPermanently, "/auth")
	}

	history, err := h.logic.History(cookie.Value)
	var t = schema.Response{Text: history}
	if err != nil {
		t = schema.Response{Text: err.Error()}
	}

	bytes, err := json.Marshal(t)
	if err != nil {
		t = schema.Response{Text: err.Error()}
		bytes, err = json.Marshal(t)
		if err != nil {
			h.Warn(err.Error())
		}
	}
	c.Response().Write(bytes)
	return nil
}

func (h *Handler) Servers(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if cookie == nil || err != nil {
		return c.Redirect(http.StatusMovedPermanently, "/auth")
	}

	servers, err := h.logic.Servers(c.Request().Context(), cookie.Value)
	if servers == "" {
		servers = "no available balancer"
	}

	var t = schema.Response{Text: servers}
	if err != nil {
		t = schema.Response{Text: err.Error()}
	}

	bytes, err := json.Marshal(t)
	if err != nil {
		t = schema.Response{Text: err.Error()}
		bytes, err = json.Marshal(t)
		if err != nil {
			h.Warn(err.Error())
		}
	}
	c.Response().Write(bytes)
	return nil
}

func (h *Handler) Authenticate(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if cookie != nil && err == nil {
		return c.Redirect(http.StatusMovedPermanently, "/")
	}

	ctx := c.Request().Context()

	username := c.FormValue("username")
	password := c.FormValue("password")

	token, err := h.logic.Authenticate(ctx, username, password)
	if err != nil && (h.security.On && h.security.MandatoryAuthorization) {
		return c.Redirect(http.StatusMovedPermanently, "/auth")
	}

	cookie = cookies.SetCookie(token)
	c.SetCookie(cookie)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func (h *Handler) Exit(c echo.Context) error {
	// TODO: invalidate token on server
	c.SetCookie(nil)
	return c.Redirect(http.StatusMovedPermanently, "/auth")
}
