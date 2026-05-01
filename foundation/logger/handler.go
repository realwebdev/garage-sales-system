package logger

import (
	"context"
	"log/slog"
)

// logHandler provides a wrapper arount the slog handler to capture which
// log level is abeing logged for event handling
type logHandler struct {
	handler slog.Handler
	events  Events
}

func newLogHandler(handler slog.Handler, events Events) *logHandler {
	return &logHandler{
		handler: handler,
		events:  events,
	}
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs returns a new JSONHandler whose attributes consists
// of h's attributes followed by attrs.
func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logHandler{
		handler: h.handler.WithAttrs(attrs),
		events:  h.events,
	}
}

// WithGroup returns a new Handler with the given group appended to the
// receiver's existing groups. The keys of all subsequent attributes, whether
// added by With or in a Record, should be qualified by the sequence of group
// names.
func (h *logHandler) WithGroup(name string) slog.Handler {
	return &logHandler{handler: h.handler.WithGroup(name), events: h.events}
}

func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	switch r.Level {
	case slog.LevelDebug:
		if h.events.Debug != nil {
			h.events.Debug(ctx, toRecord(r))
		}

	case slog.LevelError:
		if h.events.Error != nil {
			h.events.Debug(ctx, toRecord(r))
		}

	case slog.LevelWarn:
		if h.events.Warn != nil {
			h.events.Debug(ctx, toRecord(r))
		}

	case slog.LevelInfo:
		if h.events.Info != nil {
			h.events.Debug(ctx, toRecord(r))
		}
	}

	return h.handler.Handle(ctx, r)
}
