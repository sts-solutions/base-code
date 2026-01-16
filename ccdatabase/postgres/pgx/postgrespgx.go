package pgx

import (
	"context"
	"database/sql"
	"time"

	"emperror.dev/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sts-solutions/base-code/ccdatabase"
	"github.com/sts-solutions/base-code/ccdatabase/migration"
	"github.com/sts-solutions/base-code/ccdatabase/postgres"
	"github.com/sts-solutions/base-code/ccretry"
)

// postgresPgxDB implements the ccdatabase.Database interface using
// pgxpool for connection pooling and supporting migrations with sql.DB.
type postgresPgxDB struct {
	connection postgres.DBConnection // Holds DB connection details and connection string
	cfg        *pgxpool.Config       // pgxpool configuration
	migrations migration.Migration   // Migration engine to run DB migrations
	timeout    time.Duration         // Timeout applied to DB operations
	pool       *pgxpool.Pool         // pgx connection pool
	sqlDb      *sql.DB               // sql.DB instance required by Goose migrations
	closeFunc  func() error          // Function used to close all underlying resources
}

// Open initializes the pg connection pool and prepares the close function.
// It also creates a sql.DB instance because migrations require the database/sql API.
func (d *postgresPgxDB) Open() error {
	// Create the pgx pool
	pool, err := pgxpool.NewWithConfig(context.Background(), d.cfg)
	if err != nil {
		return errors.Wrap(err, "creating pgx pool")
	}
	d.pool = pool

	// Create the sql.DB wrapper required for database migrations
	sqlDb, err := sql.Open("pgx", d.connection.ConnectionString())
	if err != nil {
		return errors.Wrap(err, "creating sql DB for migrations")
	}
	d.sqlDb = sqlDb

	// Prepare a function to properly close pool and sql.DB resources
	d.closeFunc = func() error {
		var errs []error

		if d.pool != nil {
			d.pool.Close()
		}

		if d.sqlDb != nil {
			if err := d.sqlDb.Close(); err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return errors.Errorf("errors closing db: %v", errs)
		}
		return nil
	}

	return nil
}

// Close closes the database connection/pool and all its dependencies.
// Consumers should call this when the database is no longer needed.
func (d *postgresPgxDB) Close() error {
	return d.closeFunc()
}

func (d *postgresPgxDB) Conn() any {
	return d.pool
}

func (d *postgresPgxDB) Ping(ctx context.Context) error {
	if d.pool == nil {
		return errors.New("database pool is not initialized")
	}
	return d.pool.Ping(ctx)
}

func (d *postgresPgxDB) Setup(ctx context.Context) (ccdatabase.Database, error) {
	if err := d.Open(); err != nil {
		return nil, errors.Wrap(err, "opening database")
	}

	res, err := ccretry.
		NewRetry(func() error {
			return d.Ping(ctx)
		}).
		WithMaxAttempts(5).
		WithSleep(time.Second).
		Run()

	if err != nil {
		return nil, errors.Wrapf(err, "pinging database\n%s", res.String())
	}
	// Run migrations if provided
	if d.migrations != nil {
		if err := d.migrations.SetDB(d.sqlDb); err != nil {
			return nil, errors.Wrap(err, "running migrations")
		}

		if err := d.migrations.Run(ctx); err != nil {
			return nil, errors.Wrap(err, "running migrations up")
		}
	}

	return d, nil
}

func (d *postgresPgxDB) DBContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d.timeout)
}

func (d *postgresPgxDB) DBName() string {
	return d.connection.DBName
}
