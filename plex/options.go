package plex

import "log/slog"

func WithCacheLibrary(location string) func(*Plex) {
	return func(p *Plex) {
		p.cacheLibrary = location
	}
}
func WithLogger(l *slog.Logger) func(*Plex) {
	return func(p *Plex) {
		p.logger = l
	}
}
