package context_manager

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type ContextKey string

const (
	Connection     ContextKey = "connection"
	Token          ContextKey = "token"
	TransactionKey ContextKey = "transaction_key"
	UserClaim      ContextKey = "user_claim"
)

func GetTransactionContext(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(TransactionKey).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}

func GetConnectionContext(ctx context.Context) (string, error) {
	if value, ok := ctx.Value(Connection).(string); ok {
		return value, nil
	}
	return "", errors.New("connection context invalid")
}

func WithConnection(ctx context.Context, connectionString string) context.Context {
	return context.WithValue(ctx, Connection, connectionString)
}

func WithTransaction(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, TransactionKey, tx)
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

