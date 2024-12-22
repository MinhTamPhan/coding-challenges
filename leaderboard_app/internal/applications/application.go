package applications

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/wire"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
)

var (
	RestApplicationProviders = wire.NewSet(
		ProvideHttpRestful,
		ProvideHttpMetric,
		NewRestApplication,
	)
)

type metricServer *http.Server
type restServer *http.Server

func ProvideHttpRestful(httpHandler http.Handler) restServer {
	return &http.Server{
		Addr:    ":8080",
		Handler: httpHandler,
	}
}

func ProvideHttpMetric() metricServer {
	e := echo.New()
	e.GET("/metrics", echoprometheus.NewHandler())
	return &http.Server{
		Addr:        ":8081",
		Handler:     e,
		ReadTimeout: 2 * time.Second,
	}
}

type ApplicationConfig struct {
	Name            string        `mapstructure:"name"`
	Version         string        `mapstructure:"version"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type Application interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type restApplication struct {
	httpServer   *http.Server
	metricServer *http.Server
}

func NewRestApplication(httpServer restServer, metricServer metricServer) Application {
	return &restApplication{
		httpServer:   httpServer,
		metricServer: metricServer,
	}
}

func (r *restApplication) Start(ctx context.Context) error {
	log.Print("Starting application")
	g, groupCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return r.httpServer.ListenAndServe()
	})
	g.Go(func() error {
		return r.metricServer.ListenAndServe()
	})
	g.Go(func() error {
		<-groupCtx.Done()
		return r.shutdown(context.Background(), 5*time.Second)
	})
	if err := g.Wait(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (r *restApplication) shutdown(ctx context.Context, shutdownTimeout time.Duration) error {
	shutDownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()
	g, groupCtx := errgroup.WithContext(shutDownCtx)
	g.Go(func() error {
		return r.httpServer.Shutdown(groupCtx)
	})
	g.Go(func() error {
		return r.metricServer.Shutdown(groupCtx)
	})
	return g.Wait()
}

func (r *restApplication) Stop(ctx context.Context) error {
	if err := r.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	if err := r.metricServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
