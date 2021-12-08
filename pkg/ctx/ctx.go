package ctx

import (
	"context"

	"github.com/vanamelnik/gophermart/model"

	"github.com/rs/zerolog"
)

const (
	userKey   ctxKey = "user"
	loggerKey ctxKey = "logger"
)

type ctxKey string

// WithUser adds the data of authenticated user to the provided context.
func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User fetches the authenticated user data from the provided context.
func User(ctx context.Context) *model.User {
	if ctxValue := ctx.Value(userKey); ctxValue != nil {
		if user, ok := ctxValue.(*model.User); ok {
			return user
		}
	}

	return nil
}

// WithLogger applies a logger to the context provided.
func WithLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// Logger gets the logger from the context provided.
func Logger(ctx context.Context) zerolog.Logger {
	if ctxValue := ctx.Value(loggerKey); ctxValue != nil {
		if logger, ok := ctxValue.(zerolog.Logger); ok {
			return logger
		}
	}
	// TODO: I don't want to overload this function with error returning...
	// I hope that during development and testing, I can track when the application panics.
	panic("ctx.Logger: no log in context")
}
