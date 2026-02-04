package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type backend struct {
	logger *slog.Logger
	wg     sync.WaitGroup
}

func (b *backend) serve() error {
	srv := &http.Server{
		ErrorLog:          slog.NewLogLogger(b.logger.Handler(), slog.LevelError),
		Addr:              fmt.Sprintf(":%d", 4000),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       60 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		b.logger.Info("received signal", slog.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			shutdownError <- err
		}
		b.logger.Info("completing background tasks")
		b.wg.Wait()
		shutdownError <- nil
	}()
	b.logger.Info("server started")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	if err = <-shutdownError; err != nil {
		return err
	}
	b.logger.Info("server shutdown")
	return nil
}

func main() {
	b := backend{
		logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
	if err := b.serve(); err != nil {
		return
	}
}
