package repo

import (
	_ "database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Repoer interface {
	Type() string // "sqlite"
	Open() error
	Close() error
	Verify() error
	Handle() *sqlx.DB
	ForceEmpty() error // delete data but not tables
	ForceExistDBandTables() error
}

type SingleFileDBer interface {
	Path() string // sthg.db
	Repoer
	Backupable
	// NewSingleFileDBer(relLocation string) (absLocation string, err error)
}
