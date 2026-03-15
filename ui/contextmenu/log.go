package contextmenu

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type LogEntry struct {
	Index   int    `json:"index"`
	Message string `json:"message"`
}

var (
	LogEventNewMessages string = "log_event_new_messages"
)

type DebugLog struct {
	ctx            context.Context
	LatestIndex    int
	Enabled        bool
	pendingEntries []LogEntry
}

func NewDebugLog(ctx context.Context) *DebugLog {
	return &DebugLog{ctx: ctx}
}

func (a *DebugLog) GetInitialLogs() []LogEntry {
	// Clear out the pending logs. Since they'll be sent over to the client
	defer func() {
		a.pendingEntries = []LogEntry{}
		a.Enabled = true
	}()
	entries := a.pendingEntries
	a.LatestIndex = len(entries)

	return entries

}

func (a *DebugLog) NewLogs(logs []string) {
	// store the logs as pending if the debug logger hasn't been initialized by the client yet
	if !a.Enabled {
		for _, x := range logs {
			a.LatestIndex = a.LatestIndex + 1
			a.pendingEntries = append(a.pendingEntries, LogEntry{Index: a.LatestIndex, Message: x})
		}
		runtime.EventsEmit(a.ctx, LogEventNewMessages, a.pendingEntries)
	} else {
		entries := make([]LogEntry, len(logs))
		for _, x := range logs {
			a.LatestIndex = a.LatestIndex + 1
			entries = append(entries, LogEntry{Index: a.LatestIndex, Message: x})
		}
		runtime.EventsEmit(a.ctx, LogEventNewMessages, entries)
	}
}

func (a *DebugLog) NewLogger() *slog.Logger {
	handler := NewDebugLogHandler(a)
	l := slog.New(handler)
	return l
}

// Implement an slog handler that pipes logs as events to client
type debugLogHandler struct {
	debugLog *DebugLog
}

func NewDebugLogHandler(d *DebugLog) *debugLogHandler {
	return &debugLogHandler{debugLog: d}
}

func (h *debugLogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h *debugLogHandler) Handle(ctx context.Context, rec slog.Record) error {
	str := fmt.Sprintf("%v [%v] %v",
		rec.Time.Format(time.RFC822),
		rec.Level,
		rec.Message,
	)
	rec.Attrs(func(a slog.Attr) bool {
		str += fmt.Sprintf(" :: %v=%v", a.Key, a.Value)
		return true
	})
	h.debugLog.NewLogs([]string{str})
	return nil
}

func (h *debugLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &debugLogHandler{}
}

func (h *debugLogHandler) WithGroup(name string) slog.Handler {
	return &debugLogHandler{}
}
