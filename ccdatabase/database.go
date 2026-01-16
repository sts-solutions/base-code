package ccdatabase

import (
	"context"
)

type Database interface {
	Open() (err error)
	Close() error
	Conn() any
	Ping(ctx context.Context) error
	Setup(ctx context.Context) (Database, error)
	DBContext(ctx context.Context) (context.Context, context.CancelFunc)
	DBName() string
}
