package db

import (
	"chat_app_backend/internal/sqlc/db_queries"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction = func(queries *db_queries.Queries) error

type IDbConnection interface {
	CreateTransaction(ctx context.Context, transaction Transaction) error
	GetQueries() *db_queries.Queries
	Close()
}

type Connection struct {
	queries *db_queries.Queries
	pool    *pgxpool.Pool
}

func (conn *Connection) Close() {
	conn.pool.Close()
}

func (conn *Connection) GetQueries() *db_queries.Queries {
	return conn.queries
}

func (conn *Connection) CreateTransaction(ctx context.Context, transaction Transaction) error {
	tx, err := conn.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	if err := transaction(conn.queries.WithTx(tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func CreateConnection(config *PostgresConfig, ctx *context.Context) (*Connection, error) {
	connectionPool, connectionError := pgxpool.New(*ctx, config.GetConnectionString())
	if connectionError != nil {
		return nil, connectionError
	}

	return &Connection{queries: db_queries.New(connectionPool), pool: connectionPool}, nil
}
