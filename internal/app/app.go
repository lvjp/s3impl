package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/lvjp/s3impl/pkg/s3router"
	"github.com/rs/zerolog"
)

type App struct {
	ctx    context.Context
	server *http.Server
}

func New(ctx context.Context, config Config) (*App, error) {
	app := &App{
		ctx: ctx,
		server: &http.Server{
			Addr:              config.Endpoint.Addr,
			ReadHeaderTimeout: config.Endpoint.HTTPReadHeaderTimeout,
			Handler:           s3router.New(zerolog.Ctx(ctx)),
		},
	}

	return app, nil
}

func (app *App) Run() error {
	zerolog.Ctx(app.ctx).Info().Msg("app: Start to listen and serve")

	err := app.server.ListenAndServe()
	switch {
	case errors.Is(err, http.ErrServerClosed):
		return nil
	case err != nil:
		return fmt.Errorf("app: listen error: %w", err)
	default:
		return nil
	}
}

func (app *App) Shutdown(ctx context.Context) error {
	zerolog.Ctx(app.ctx).Info().Msg("app: Shutdown")

	if err := app.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("app: shutdown error: %w", err)
	}

	return nil
}
