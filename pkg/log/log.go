package log

import (
	"context"

	"github.com/Sirupsen/logrus"
)

func WithContext(ctx context.Context) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"request-id":      ctx.Value(IDKey),
		"request-ip":      ctx.Value(IPKey),
		"request-user-id": ctx.Value(UserIDKey),
	})
}
