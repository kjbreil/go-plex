package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kjbreil/go-plex/pkg/plex"
)

const webhookPort = 8081

func main() {
	plexHost := os.Getenv("PLEX_HOST")
	plexToken := os.Getenv("PLEX_TOKEN")
	logger := slog.Default()
	conn, err := plex.New(
		plexHost,
		plexToken,
		plex.WithCacheLibrary("plex-library-cache.json"),
		plex.WithLogger(logger),
	)
	if err != nil {
		panic(err)
	}

	err = conn.InitLibraries()
	if err != nil {
		panic(err)
	}

	go func() {
		start := time.Now()
		conn.PopulateLibraries()()
		logger.Info("plex library refreshed", "duration", time.Since(start))
	}()

	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = setupWebhooks(conn)
	}()

	<-ctrlC

	conn.Close()
}

func setupWebhooks(conn *plex.Plex) error {
	ip := net.ParseIP("10.0.2.2")
	logger := slog.Default()

	conn.Webhook = plex.NewWebhook(webhookPort, ip)

	if err := conn.Webhook.OnPlay(func(w plex.WebhookEvent) {
		logger.Info("media playing", "title", w.Metadata.Title)
	}); err != nil {
		return err
	}

	if err := conn.Webhook.OnPause(func(w plex.WebhookEvent) {
		logger.Info("media paused", "title", w.Metadata.Title)
	}); err != nil {
		return err
	}

	if err := conn.Webhook.OnResume(func(w plex.WebhookEvent) {
		logger.Info("media resumed", "title", w.Metadata.Title)
	}); err != nil {
		return err
	}

	if err := conn.Webhook.OnStop(func(w plex.WebhookEvent) {
		logger.Info("media stopped", "title", w.Metadata.Title)
	}); err != nil {
		return err
	}

	if err := conn.Webhook.OnRate(func(w plex.WebhookEvent) {
		logger.Info("media rated", "title", w.Metadata.Title)
	}); err != nil {
		return err
	}

	if err := conn.Webhook.OnScrobble(func(w plex.WebhookEvent) {
		logger.Info("media scrobbled", "title", w.Metadata.Title)
	}); err != nil {
		return err
	}

	conn.ServeWebhook()

	return nil
}
