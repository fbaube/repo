package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	S "strings"

	// DU "github.com/fbaube/dbutils"
	L "github.com/fbaube/mlog"
)

// CreateTable_sqlite creates a table for our simplified SQLite DB model where
// - Every column is either string ("TEXT") or int ("INTEGER"),
// - Every column is NOT NULL (because NULL is evil),
// - Every column has type checking (TBS), and
// - Every table has a primary index field, and
// - Every index (both primary and foreign) includes the full name of the
// table, which simplifies column creation and table cross-referencing
// (and in particular, JOINs).
//
func (pDB *SimpleRepo) CreateTable_sqlite(ts TableConfig) error {
	if pDB.Type() != "sqlite" {
		return errors.New("simplerepo.sqlite.createtable: not sqlite")
	}
	var CTS string // the Create Table SQL string
	var hasFKs bool
	hasFKs = (ts.ForenKeys != nil && len(ts.ForenKeys) > 0)

	// === CREATE TABLE
	CTS = "CREATE TABLE " + ts.TableName + "(\n"
	// == PRIMARY KEY
	CTS += "idx_" + ts.TableName + " integer not null primary key autoincrement, "
	CTS += "-- NOTE: integer, not int. \n"
	if hasFKs {
		// === FOREIGN KEYS
		// []string{"map_contentity", "tpc_contentity"},
		for _, tbl := range ts.ForenKeys {
			if S.Contains(tbl, "_") {
				i := S.LastIndex(tbl, "_")
				minTbl := tbl[i+1:]
				L.L.Info("DB compound index: " + tbl + " indexes " + minTbl)
				CTS += "idx_" + tbl + " integer not null references " + minTbl + ", \n"
			} else {
				// idx_inb integer not null references INB,
				// "not null" might be problematic during development.
				CTS += "idx_" + tbl + " integer not null references " + tbl + ", \n"
			}
		}
	}
	for _, fld := range ts.Columns {
		switch fld.TxtIntKeyEtc {
		case D_INT:
			// e.g.: filect int not null check (filect >= 0) default 0
			// also: `Col1 INTEGER CHECK (typeof(Col1) == 'integer')`
			//
			CTS += fld.Code + " int not null"
			// CTS += fld.Code + " int not null check (typeof(" + fld.Code + ") == 'int')"
			/*
				switch ts.intRanges[i] {
				case 1:
					// check (filect >= 0)
					CTS += " check (" + fld + " > 0), \n"
				case 0:
					CTS += " check (" + fld + " >= 0), \n"
				default: // case -1:
					CTS += ", \n"
				}
			*/
			CTS += ", \n"
		case D_TXT:
			CTS += fld.Code + " text not null check (typeof(" + fld.Code + ") == 'text'), \n"
		default:
			panic("Unhandled: " + fld.TxtIntKeyEtc)
		}
	}
	if hasFKs {
		// FOREIGN KEY(idx_inb) REFERENCES INB(idx_inb)
		for _, tbl := range ts.ForenKeys {
			// idx_inb integer not null references INB,
			// TMP := "foreign key(idx_" + tbl + ") references " + tbl + "(idx_" + tbl + "), \n"
			// println("TMP:", TMP)
			CTS += "foreign key(idx_" + tbl + ") references " + tbl + "(idx_" + tbl + "), \n"
		}
	}

	CTS = S.TrimSuffix(CTS, "\n")
	CTS = S.TrimSuffix(CTS, " ")
	CTS = S.TrimSuffix(CTS, ",")
	CTS += "\n);"

	fnam := "./create-table-" + ts.TableName + ".sql"
	e := ioutil.WriteFile(fnam, []byte(CTS), 0644)
	if e != nil {
		L.L.Error("Could not write file: " + fnam)
	} else {
		L.L.Dbg("Wrote \"CREATE TABLE " + ts.TableName + " ... \" to: " + fnam)
	}
	pDB.Handle().MustExec(CTS)
	ss, e := pDB.DumpTableSchema_sqlite(ts.TableName)
	if e != nil {
		return fmt.Errorf("simplerepo.createtable.sqlite: "+
			"DumpTableSchema<%s> failed: %w", e)
	}
	L.L.Dbg(ts.TableName + " SCHEMA: " + ss)
	// println("TODO: Insert record with IDX 0 and string descr's")
	//    and ("TODO: Dump all table records (i.e. just one)")
	return nil
}

// DbTblColsInDb returns all column names & types for the specified table.
//
func (pDB *SimpleRepo) DbTblColsInDb(tableName string) ([]*DbColInDb, error) {
	if tableName == "" {
		return nil, nil
	}
	if pDB.Type() != "sqlite" {
		return nil, errors.New(
			"simplerepo.sqlite.dbtblcolsindb: not sqlite")
	}
	var e error
	var Rs *sql.Rows
	var CTs []*sql.ColumnType
	var retval []*DbColInDb

	Rs, e = pDB.Handle().Query("select * from " + tableName + " limit 1")
	if e != nil {
		return nil, fmt.Errorf("simplerepo.dbtblcolsindb: "+
			"select * : failed on table <%s>: %w", tableName, e)
	}
	CTs, e = Rs.ColumnTypes()
	if e != nil {
		return nil, fmt.Errorf("simplerepo.dbtblcolsindb: "+
			"rs.ColumnTypes failed on table <%s>: %w", tableName, e)
	}
	for _, ct := range CTs {
		dci := new(DbColInDb)
		dci.TxtIntKeyEtc = TxtIntKeyEtc(ct.DatabaseTypeName())
		dci.Code = ct.Name()
		retval = append(retval, dci)
	}
	return retval, nil
}

// DumpTableSchema_sqlite returns the SQLite schema for the specified table.
//
func (pDB *SimpleRepo) DumpTableSchema_sqlite(tableName string) (string, error) {
	if pDB.Type() != "sqlite" {
		return "", errors.New(
			"simplerepo.sqlite.dumptableschema: not sqlite")
	}
	var theCols []*DbColInDb
	var sb S.Builder
	var sType string
	var e error

	theCols, e = pDB.DbTblColsInDb(tableName)
	if e != nil {
		return "", fmt.Errorf("simplerepo.dumptableschema.sqlite: "+
			"pDB.DbTblColsInDb<%s> failed: %w", e)
	}
	for i, c := range theCols {
		sType = ""
		if c.TxtIntKeyEtc != "text" {
			sType = "(" + string(c.TxtIntKeyEtc) + "!)"
		}
		sb.Write([]byte(fmt.Sprintf("[%d]%s%s / ", i, sType, c.Code)))
	}
	sb.Write([]byte(fmt.Sprintf("%d fields", len(theCols))))
	return sb.String(), nil
}
