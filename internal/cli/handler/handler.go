package handler

import (
	"fmt"
	"github.com/egorgasay/grpc-storage/internal/cli/config"
	"github.com/egorgasay/grpc-storage/internal/cli/usecase"
	"github.com/egorgasay/grpc-storage/pkg/logger"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/renderer"
	"net/http"
)

var rnd *renderer.Render

func init() {
	opts := renderer.Options{
		ParseGlobPattern: "templates/html/*.html",
	}

	rnd = renderer.New(opts)
}

type Handler struct {
	cfg   *config.Config
	logic *usecase.UseCase
	logger.ILogger
}

func New(cfg *config.Config, logic *usecase.UseCase, loggerInstance logger.ILogger) *Handler {
	return &Handler{cfg: cfg, logic: logic, ILogger: loggerInstance}
}

func (h *Handler) MainPage(e echo.Context) error {
	err := e.Render(http.StatusOK, "index.html", nil)
	if err != nil {
		h.Warn(err.Error())
	}
	return err
}

func (h *Handler) Action(e echo.Context) error {
	e.Response().Header().Set("Content-Type", "application/json")
	action := e.Request().URL.Query().Get("action")
	res, err := h.logic.ProcessQuery(h.cfg, action)
	if err != nil {
		h.Info(err.Error())
		e.Response().Write([]byte(fmt.Sprintf(`{"text": "%v"}`, err)))
		e.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	resp := fmt.Sprintf(`{"text": "%v"}`, res)
	e.Response().Write([]byte(resp))
	e.Response().WriteHeader(http.StatusOK)
	return nil
}

//func (h *Handler) Set() {
//
//}
//
//func (h *Handler) Get() () {
//
//}
