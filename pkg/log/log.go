package log

import (
	"context"

	"github.com/Sirupsen/logrus"
)

var logrusger = logrus.New()

func WithContext(ctx context.Context) *logrus.Entry {
	return logrusger.WithField("request-id", ctx.Value(IDKey))
}
