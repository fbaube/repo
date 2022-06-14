package repo

import (
	_ "database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

/*
type Repoer interface {
        Type() string // "sqlite"
        Open() error
        Close() error
        Verify() error
        Handle() *sqlx.DB
        ForceEmpty() error // delete data but not tables
        ForceExistDBandTables() error
}
type Backupable interface {
        MoveToBackup() (string, error)
        CopyToBackup() (string, error)
        RestoreFromMostRecentBackup() (string, error)
}*/

// SimpleRepo can be fully described by (1) a filepath or a URL,
// plus (2) a "type" string (tipicly "sqlite"). A SimpleRepo is
// expected to be Backupable, and most likely a sSingleFileDBer.
//
type SimpleRepo struct {
	db        *sqlx.DB
	pathOrUrl string
	isPath    bool
	isUrl     bool
	tipe      string // "sqlite"
}

func (p *SimpleRepo) Type() string {
	return p.tipe
}

func (p *SimpleRepo) Path() string {
	if !p.isPath {
		return ""
	}
	return p.pathOrUrl
}

func (p *SimpleRepo) Handle() *sqlx.DB {
	return p.db
}

// setHandle is limited to package scope.
func (p *SimpleRepo) setHandle(pDB *sqlx.DB) {
	p.db = pDB
}
