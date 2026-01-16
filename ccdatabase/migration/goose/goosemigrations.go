package goose

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sts-solutions/base-code/ccdatabase/migration"
	"github.com/sts-solutions/base-code/cclog"

	"github.com/pressly/goose"
)

type gooseMigrations struct {
	sourceFolder string
	logger       gooseLogger
	db           *sql.DB
}

func NewGooseMigrations(sourceFolder string,
	logger cclog.Logger) migration.Migration {
	return &gooseMigrations{
		sourceFolder: sourceFolder,
		logger:       newGooseLogger(logger),
	}
}
func (g *gooseMigrations) Run(ctx context.Context) error {
	loggerContext := g.logger.ctx
	defer func() {
		g.logger.ctx = loggerContext
	}()

	for {
		_, err := g.db.ExecContext(ctx, "CREATE TABLE goose_migrations_in_progress (dummy boolean)")
		if err != nil {
			break
		}

		g.logger.Print("cannot acquire migration lock, another migration in progress: %s", err)
		time.Sleep(time.Second)
	}

	defer func() {
		_, err := g.db.ExecContext(ctx, "DROP TABLE IF EXISTS goose_migrations_in_progress")
		if err != nil {
			g.logger.Print("failed to drop migration table: %s", err)
		}
	}()

	goose.SetLogger(g.logger)

	return goose.Run("Up", g.db, g.sourceFolder)
}

func (g *gooseMigrations) SetDB(db any) error {
	sqlDb, ok := db.(*sql.DB)
	if !ok {
		return fmt.Errorf("invalid db type: expected *sql.DB, got %T, expected *sql.DB", db)
	}
	g.db = sqlDb
	return nil
}

func (g *gooseMigrations) SourceFolder() string {
	return g.sourceFolder
}
