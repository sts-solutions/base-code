package pgx

import (
	"database/sql"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sts-solutions/base-code/ccdatabase"
	"github.com/sts-solutions/base-code/ccdatabase/postgres"
	"github.com/sts-solutions/base-code/ccerrors"
	"github.com/sts-solutions/base-code/ccvalidation"
)

type postgresPgxDBBuilder struct {
	postgresPgxDB         *postgresPgxDB
	maxConns              int
	minIdleConns          int
	connMaxLifeTime       time.Duration
	connMaxLifeTimeJitter time.Duration
	connMaxIdleTime       time.Duration
	healthCheckInterval   time.Duration
}

func NewPostgresPgxDBBuilder() *postgresPgxDBBuilder {
	return &postgresPgxDBBuilder{
		postgresPgxDB: &postgresPgxDB{
			connection: postgres.DBConnection{},
			sqlDb:      &sql.DB{},
			closeFunc: func() error {
				return ccerrors.ErrNotImplemented
			},
		},
	}
}

func (b *postgresPgxDBBuilder) Build() (ccdatabase.Database, error) {
	err := b.validate()
	if err != nil {
		return nil, err
	}

	b.postgresPgxDB.cfg, err = pgxpool.ParseConfig(b.postgresPgxDB.connection.ConnectionString())
	if err != nil {
		return nil, errors.Wrap(err, "parsing pgx pool config")
	}

	if b.maxConns > 0 {
		b.postgresPgxDB.cfg.MaxConns = int32(b.maxConns)
	}

	if b.minIdleConns > 0 {
		b.postgresPgxDB.cfg.MinConns = int32(b.minIdleConns)
	}

	if b.connMaxLifeTime > 0 {
		b.postgresPgxDB.cfg.MaxConnLifetime = b.connMaxLifeTime
	}

	if b.connMaxLifeTimeJitter > 0 {
		b.postgresPgxDB.cfg.MaxConnLifetimeJitter = b.connMaxLifeTimeJitter
	}

	if b.connMaxIdleTime > 0 {
		b.postgresPgxDB.cfg.MaxConnIdleTime = b.connMaxIdleTime
	}

	if b.healthCheckInterval > 0 {
		b.postgresPgxDB.cfg.HealthCheckPeriod = b.healthCheckInterval
	}

	return b.postgresPgxDB, nil

}

func (b *postgresPgxDBBuilder) WithTimeout(timeout time.Duration) *postgresPgxDBBuilder {
	b.postgresPgxDB.timeout = timeout
	return b
}

func (b *postgresPgxDBBuilder) WithMinIdleConns(minIdleConns int) *postgresPgxDBBuilder {
	b.minIdleConns = minIdleConns
	return b
}

func (b *postgresPgxDBBuilder) WithConnMaxLifetime(connMaxLifetime time.Duration) *postgresPgxDBBuilder {
	b.connMaxLifeTime = connMaxLifetime
	return b
}

func (b *postgresPgxDBBuilder) WithConnMaxLifetimeJitter(connMaxLifetimeJitter time.Duration) *postgresPgxDBBuilder {
	b.connMaxLifeTimeJitter = connMaxLifetimeJitter
	return b
}

func (b *postgresPgxDBBuilder) WithConnMaxIdleTime(connMaxIdleTime time.Duration) *postgresPgxDBBuilder {
	b.connMaxIdleTime = connMaxIdleTime
	return b
}

func (b *postgresPgxDBBuilder) WithHealthCheckInterval(healthCheckInterval time.Duration) *postgresPgxDBBuilder {
	b.healthCheckInterval = healthCheckInterval
	return b
}

func (b *postgresPgxDBBuilder) WithConnectionString(connectionString string) *postgresPgxDBBuilder {
	b.postgresPgxDB.connection = postgres.NewDBConnectionFromCnnStringUrl(connectionString)
	return b
}

func (b *postgresPgxDBBuilder) WithConnectionParams(p postgres.DBConnection) *postgresPgxDBBuilder {
	b.postgresPgxDB.connection = *postgres.NewDBConnectionParams(p.Host,
		p.DBName,
		p.SSLMode,
		p.Port,
		p.UserName,
		p.Password)
	return b
}

func (b *postgresPgxDBBuilder) validate() error {
	res := ccvalidation.Result{}

	if b.postgresPgxDB.timeout <= 0 {
		res.AddErrorMessage(
			fmt.Sprintf("timeout must be greater than zero, got %v", b.postgresPgxDB.timeout),
		)
	}

	if err := b.postgresPgxDB.connection.Validate(); err != nil {
		for _, e := range err.(ccvalidation.Result).GetFailures() {
			res.AddFailure(e)
		}
	}

	if res.IsFailure() {
		return res
	}

	return nil
}
