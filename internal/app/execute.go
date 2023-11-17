package app

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc/pool"
	"gopkg.in/yaml.v3"
)

func Execute(configPath string) error {
	config, err := readConfiguration(configPath)
	if err != nil {
		return fmt.Errorf("could not read config: %w", err)
	}

	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger.With().Str("module", "stdlog").Logger())

	app, err := New(log.Logger.WithContext(context.Background()), *config)
	if err != nil {
		return fmt.Errorf("could not initialize the application: %w", err)
	}

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool := pool.New().WithContext(mainCtx)

	pool.Go(func(_ context.Context) error {
		if err := app.Run(); err != nil {
			return fmt.Errorf("could not run app: %w", err)
		}

		return nil
	})

	pool.Go(func(ctx context.Context) error {
		<-ctx.Done()

		if err := app.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("could not shutdown app: %w", err)
		}

		return nil
	})

	<-mainCtx.Done()
	log.Info().Msg("Shutdown started")
	if err := pool.Wait(); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}
	log.Info().Msg("Shutdown finished")

	return nil
}

func readConfiguration(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open failed: %w", err)
	}
	defer file.Close()

	var config Config
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	return &config, nil
}
