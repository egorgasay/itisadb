package handler

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/internal/cli/config"
	"itisadb/internal/cli/cookies"
	"itisadb/internal/cli/schema"
	"itisadb/internal/cli/usecase"
	"itisadb/pkg/logger"
	"net/http"
)

type Handler struct {
	cfg   *config.Config
	logic *usecase.UseCase
	logger.ILogger
}

func New(cfg *config.Config, logic *usecase.UseCase, loggerInstance logger.ILogger) *Handler {
	return &Handler{cfg: cfg, logic: logic, ILogger: loggerInstance}
}

func (h *Handler) MainPage(c echo.Context) error {
	err := c.Render(http.StatusOK, "index.html", nil)
	if err != nil {
		h.Warn(err.Error())
	}
	return err
}

func (h *Handler) Action(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	cookie, err := c.Cookie("session")
	if err != nil || cookie == nil {
		cookie = cookies.SetCookie()
		c.SetCookie(cookie)
	}

	action := c.Request().URL.Query().Get("action")
	res, err := h.logic.ProcessQuery(cookie.Value, action)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code().String() == codes.NotFound.String() {
			err = errors.New("the value is not found")
		} else if st.Message() == "unknown server" {
			err = errors.New("unknown server")
		} else if ok && st.Code().String() == codes.Unavailable.String() {
			err = errors.New("memory balancer is offline")
		}

		h.Warn(err.Error())
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
		cookie = cookies.SetCookie()
		c.SetCookie(cookie)
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
	servers, err := h.logic.Servers()
	if servers == "" {
		servers = "no available servers"
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
