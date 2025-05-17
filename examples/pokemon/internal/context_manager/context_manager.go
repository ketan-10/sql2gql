package context_manager

import (
	"context"
	"errors"

)

type ContextKey string

const (
	Connection     ContextKey = "connection"
	Token          ContextKey = "token"
	UserClaim      ContextKey = "user_claim"
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


func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, Token, token)
}

func GetTokenContext(ctx context.Context) (string, error) {
	if value, ok := ctx.Value(Token).(string); ok {
		return value, nil
	}
	return "", errors.New("token does not extis")
}

