package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kjbreil/go-plex/pkg/plex"
)

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

	go setupWebhooks(conn)

	<-ctrlC

	conn.Close()
}

func setupWebhooks(conn *plex.Plex) error {
	ip := net.ParseIP("10.0.2.2")

	conn.Webhook = plex.NewWebhook(8081, ip)

	conn.Webhook.OnPlay(func(w plex.WebhookEvent) {
		fmt.Printf("%s is playing\n", w.Metadata.Title)
	})

	conn.Webhook.OnPause(func(w plex.WebhookEvent) {
		fmt.Printf("%s is paused\n", w.Metadata.Title)
	})

	conn.Webhook.OnResume(func(w plex.WebhookEvent) {
		fmt.Printf("%s has resumed\n", w.Metadata.Title)
	})

	conn.Webhook.OnStop(func(w plex.WebhookEvent) {
		fmt.Printf("%s has stopped\n", w.Metadata.Title)
	})
	conn.Webhook.OnRate(func(w plex.WebhookEvent) {
		fmt.Printf("%s has been rated\n", w.Metadata.Title)
	})

	conn.Webhook.OnScrobble(func(w plex.WebhookEvent) {
		fmt.Printf("%s has been scrobbled\n", w.Metadata.Title)
	})

	conn.ServeWebhook()

	return nil
}
