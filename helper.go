package logrus

import (
	"context"

	"github.com/sirupsen/logrus"
)

// EntryKey is the context key for the logrus.Entry.
//
// It use a private struct in order to avoid any conflict with other packages.
// Its value of zero is arbitrary.
const EntryKey key = 0

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// LogError add an error into the context.
//
// This error will be added as a field inside the request log.
func LogError(ctx context.Context, err error) {
	entry := ctx.Value(EntryKey).(*logrus.Entry)

	entry.Data[logrus.ErrorKey] = err.Error()
}
