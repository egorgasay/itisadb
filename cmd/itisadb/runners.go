package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"

	"github.com/brpaz/echozap"
	"github.com/egorgasay/gost"
	api "github.com/egorgasay/itisadb-shared-proto/go"
	"github.com/go-chi/chi/middleware"
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"itisadb/config"
	"itisadb/internal/cli/handler"
	"itisadb/internal/cli/storage"
	"itisadb/internal/cli/usecase"
	"itisadb/internal/domains"
	grpchandler "itisadb/internal/handler/grpc"
	resthandler "itisadb/internal/handler/rest"
	"itisadb/internal/service/balancer"
)

func runGRPC(
	ctx context.Context,
	l *zap.Logger,
	logic domains.Core,
	securityCFG config.SecurityConfig,
	networkCFG config.NetworkConfig,
	session domains.Session,
) {
	h := grpchandler.New(logic, l, session, securityCFG)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(h.AuthMiddleware),
	)

	lis, err := net.Listen("tcp", networkCFG.GRPC)
	if err != nil {
		l.Fatal("failed to listen: %v", zap.Error(err))
	}
	api.RegisterItisaDBServer(grpcServer, h)

	err = gost.WithContextPool(ctx, func() error {
		l.Info("Starting GRPC", zap.String("address", networkCFG.GRPC))
		err = grpcServer.Serve(lis)
		if err != nil {
			return fmt.Errorf("grpcServer Serve: %v", err)
		}

		return nil
	}, make(chan struct{}, 1), func() {
		grpcServer.GracefulStop()
	})
	l.Info("Shutdown GRPC ...")

	if err != nil && !errors.Is(err, context.Canceled) {
		l.Fatal("Error in GRPC: %v", zap.Error(err))
	}
}

func runREST(
	ctx context.Context,
	l *zap.Logger,
	logic *balancer.Balancer,
	cfg config.NetworkConfig,
) {
	h := resthandler.New(logic)
	lis, err := net.Listen("tcp", cfg.REST)
	if err != nil {
		l.Fatal("failed to listen: %v", zap.Error(err))
	}

	err = gost.WithContextPool(ctx, func() error {
		l.Info("Starting FastHTTP %s", zap.String("address", cfg.REST))
		if err := fasthttp.Serve(lis, h.ServeHTTP); err != nil {
			return fmt.Errorf("error in REST Serve: %w", err)
		}

		return nil
	}, make(chan struct{}, 1), func() {
		if err := lis.Close(); err != nil {
			l.Warn("Failed to close listener", zap.Error(err))
		}
	})
	l.Info("Shutdown REST ...")

	if err != nil && !errors.Is(err, context.Canceled) {
		l.Fatal("Error in REST: %v", zap.Error(err))
	}
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func runWebCLI(ctx context.Context, cfg config.WebAppConfig, securityConfig config.SecurityConfig, l *zap.Logger, balancer string) {
	store := storage.New()
	logic := usecase.New(cfg, store, balancer, l)

	e := echo.New()

	if cfg.Logs {
		e.Use(echozap.ZapLogger(l))
	}

	e.Use(echo.WrapMiddleware(middleware.Recoverer))
	h := handler.New(logic, l, securityConfig)
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/html/*.html")),
	}
	e.Renderer = t
	h.PublicRoutes(e)
	e.Static("/static", "static")

	//router.Use(gzip.Gzip(gzip.BestSpeed))

	lis, err := net.Listen("tcp", cfg.Host)
	if err != nil {
		l.Fatal("failed to listen: %v", zap.Error(err))
		return
	}

	err = gost.WithContextPool(ctx, func() error {
		l.Info("Starting CLI HTTP server", zap.String("address", cfg.Host))
		err = http.Serve(lis, e)
		if err != nil {
			return fmt.Errorf("http.Serve: %w", err)
		}

		return nil
	}, make(chan struct{}, 1), func() {
		l.Info("Stopping CLI HTTP server")
		lis.Close()
	})
	if err != nil {
		return
	}
}
