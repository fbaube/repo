package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	FP "path/filepath"
	// S "strings"

	// DU "github.com/fbaube/dbutils"
	FU "github.com/fbaube/fileutils"
	// L "github.com/fbaube/mlog"
	// "github.com/fbaube/repo"
	SU "github.com/fbaube/stringutils"
	WU "github.com/fbaube/wasmutils"
	"github.com/jmoiron/sqlx"

	// to get init()
	_ "github.com/mattn/go-sqlite3"
)

/*
// SqliteRepo stores DB filepaths, DB cnxns, DB txns.
type SqliteRepo struct {
	FU.PathProps
	// Connection
	*sqlx.DB
	// Session-level open Txn (if non-nil). Relevant API:
	// func (db *sqlx.DB)    Beginx() (*sqlx.Tx, error)
	// func (db *sqlx.DB) MustBegin()  *sqlx.Tx
	// func (tx *sql.Tx)     Commit()   error
	// func (tx *sql.Tx)   Rollback()   error
	*sqlx.Tx
}
*/

// NewSimpleRepo does not open a DB, it merely checks that the given
// DIRECTORY (not file!) path is OK. That is to say, it initializes
// path-related variables but does not do more.
//
// Currently, aType must be "sqlite".
// aPath can be a relative path passed to the CLI; if it is "",
// the DB path is set to the CWD (current working directory).
//
func NewSimpleRepo(aType, aPath string) (*SimpleRepo, error) {
	if aType != "sqlite" {
		return nil, errors.New("simplerepo.new: not sqlite")
	}
	var e error
	var relFP = aPath
	if aPath == "" {
		// return nil, errors.New("newsimplerepo: no path")
		// If the DB name was unnecessarily provided,
		// trim it off to prevent problems.
		// FIXME: relFP = S.TrimSuffix(aPath, DBNAME)
		relFP, e = os.Getwd() // "."
		if e != nil {
			// TRY: https://ian-says.com/articles/golang-in-the-browser-with-web-assembly/
			// "syscall/js"
			// js.Global != nil && js.Global.Get("document") != js.Undefined
			if WU.IsBrowser() { // WU.IsWasm() {
				println("FIXME: Where is DB in browser WASM ?")
			}
			return nil, fmt.Errorf(
				"newsimplerepo: can't get CWD: %w", e)
			// os.Exit(1)
		}
	}
	pDB := new(SimpleRepo)
	var pPP *FU.PathProps
	pPP, e = FU.NewPathProps(aPath)
	if e != nil || pPP == nil {
		return nil, FU.WrapAsPathPropsError(
			e, "newsimplerepo.NewPathProps(1)", pPP)
	}
	if !pPP.IsOkayDir() {
		return nil, errors.New("DB dir not exist or not a dir: " + aPath)
	}
	pPP, e = FU.NewPathProps(FP.Join(relFP, "mmmc.db"))
	if e != nil {
		return nil, FU.WrapAsPathPropsError(
			e, "newsimplerepo.NewPathProps(2)", pPP)
	}
	// pDB.PathProps = *pPP
	return pDB, nil
}

// ForceExistDBandTables creates a new empty DB with the proper schema.
func (p *SimpleRepo) ForceExistDBandTables() error {
	if p.Path() == "" {
		return errors.New("simplerepo.forceexistdbandtables: no path")
	}
	if p.Type() != "sqlite" {
		return errors.New("simplerepo.forceexistdbandtables: not sqlite")
	}
	if p.Handle() == nil {
		return errors.New("simplerepo.forceexistdbandtables: not open")
	}
	var dest string = p.Path()
	var e error
	var theSqlDB *sql.DB

	theSqlDB, e = sql.Open("sqlite3", dest)
	if e != nil {
		return errors.New("ForceExistDBandTables: sql.Open")
	}
	e = theSqlDB.Ping()
	if e != nil {
		return errors.New("ForceExistDBandTables: sql.Ping")
	}
	e = theSqlDB.PingContext(context.Background())
	if e != nil {
		return errors.New("ForceExistDBandTables: sql.PingContext")
	}
	println("New DB created at:", SU.Tildotted(dest))
	drivers := sql.Drivers()
	println("DB driver(s):", fmt.Sprintf("%+v", drivers))
	p.setHandle(sqlx.NewDb(theSqlDB, "sqlite3"))

	for _, cfg := range AllTableConfigs {
		p.CreateTable_sqlite(cfg)
	}
	// It may seem odd that this is necessary,
	// but for some retro compatibility, SQLite does
	// not by default enforce foreign key constraints.
	p.MustExecStmt("PRAGMA foreign_keys = ON;")
	return nil
}

func (p *SimpleRepo) Verify() error {
	if p.Type() != "sqlite" {
		return errors.New("simplerepo.verify: not sqlite")
	}
	p.MustExecStmt("PRAGMA integrity_check;")
	p.MustExecStmt("PRAGMA foreign_key_check;")
	return nil
}

// ForceEmpty is a convenience function. It first makes a backup.
func (p *SimpleRepo) ForceEmpty() error {
	if p.Path() == "" {
		return errors.New("simplerepo.forcempty: no path")
	}
	if p.Type() != "sqlite" {
		return errors.New("simplerepo.forcempty: not sqlite")
	}
	if p.Handle() == nil {
		return errors.New("simplerepo.forcempty: not open")
	}
	p.MoveToBackup()
	p.ForceExistDBandTables()
	return nil
}
