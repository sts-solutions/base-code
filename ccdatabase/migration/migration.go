package migration

import (
	"context"
)

type Migration interface {
	Run(ctx context.Context) error
	SetDB(db any) error
	SourceFolder() string
}
