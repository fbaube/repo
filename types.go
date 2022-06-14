package repo

// DbDescr is generic.
type DbDescr Datum

// DbTblSpec specifies a DB table (but not its columns!).
// Field usage is as follows:
//  - TxtIntKeyEtc: D_TBL
//  - Code: a short name for use in field names (primary & foreign keys)
//  - Name: a long name, for the name of the table in the DB
//  - Descr: long description
//
type DbTblSpec DbDescr

// DbColSpec specifies a datum (i.e. a struct field and/or a
// DB column), including its generic/portable/DB-independent
// representation using the enumeration Datum.TxtIntKeyEtc).
// Some values for common DB columns are defined in the D_*
// series below. Field usage is as follows:
//  - TxtIntKeyEtc: D_TXT or D_INT
//  - Code: the field name in the DB
//  - Name: short description
//  - Descr: long description
//
type DbColSpec DbDescr

// DbColInDb describes a column as-is in the DB (as obtained via
// reflection), and has a slot to include the value (as a string).
// Not used so far.
type DbColInDb DbDescr

var D_RelFP = DbColSpec{D_TXT, "relfp", "Rel. path", "Rel.FP (from CLI)"}
var D_AbsFP = DbColSpec{D_TXT, "absfp", "Abs. path", "Absolute filepath"}
var D_TmCre = DbColSpec{D_TXT, "t_cre", "Cre. time", "Creation date+time"}
var D_TmImp = DbColSpec{D_TXT, "t_imp", "Imp. time", "DB import date+time"}
var D_TmEdt = DbColSpec{D_TXT, "t_edt", "Edit time", "Last edit date+time"}
