package context_manager

import (
	"context"
	"errors"

)

type ContextKey string

const (
	Connection     ContextKey = "connection"
)


func GetConnectionContext(ctx context.Context) (string, error) {
	if value, ok := ctx.Value(Connection).(string); ok {
		return value, nil
	}
	return "", errors.New("connection context invalid")
}

func WithConnection(ctx context.Context, connectionString string) context.Context {
	return context.WithValue(ctx, Connection, connectionString)
}


