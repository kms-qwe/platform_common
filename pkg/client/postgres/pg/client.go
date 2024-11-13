package pgv1

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	pgClient "github.com/kms-qwe/platform_common/pkg/client/postgres"
)

type client struct {
	masterDBC pgClient.DB
}

// NewPgClient create new pg client
func NewPgClient(ctx context.Context, dsn string) (pgClient.Client, error) {
	dbc, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return &client{
		masterDBC: NewDB(dbc),
	}, nil
}

// DB return DB interface
func (c *client) DB() pgClient.DB {
	return c.masterDBC
}

// Close pg client
func (c *client) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
