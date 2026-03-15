package contextmenu

import (
	"context"

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
	ctx         context.Context
	LatestIndex int
}

func NewDebugLog(ctx context.Context) *DebugLog {
	return &DebugLog{ctx: ctx}
}

func (a *DebugLog) GetInitialLogs() []LogEntry {
	entries := make([]LogEntry, 20)
	a.LatestIndex = len(entries)
	for i := range entries {
		entries[i] = LogEntry{
			Index:   i,
			Message: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		}
	}

	return entries

}

func (a *DebugLog) NewLogs(logs []string) {
	entries := make([]LogEntry, len(logs))
	for _, x := range logs {
		a.LatestIndex = a.LatestIndex + 1
		entries = append(entries, LogEntry{Index: a.LatestIndex, Message: x})
	}
	runtime.EventsEmit(a.ctx, LogEventNewMessages, entries)
}
