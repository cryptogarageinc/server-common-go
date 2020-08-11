package log

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
)

// Save adds the given fields to the logger extracted from the given context
// as a structured log entry, and returns the updated context.
func Save(ctx context.Context, fields logrus.Fields) (context.Context, *logrus.Entry) {
	entry := ctxlogrus.Extract(ctx).WithFields(fields)
	return ctxlogrus.ToContext(ctx, entry), entry
}
