package repo

// TableConfig describes the field structure of a database table,
// assuming that it is enough to use just two column types, TEXT
// (for strings) and INTEGER (for integers). Also a primary key
// is assumed and foreign keys are allowed.
//
// Note that field Columns is a slice of [DbColSpec], each of
// which is four text fields: [TxtIntKeyEtc], Code, Name, Descr.
//
// Any field can be nil or length [0], except the first (TableName).
// This means the field ForenKeys and Columns.
//
// Date-time's are not an issue for SQLite, since either a string or
// an int can be used. We will favor using strings ("TEXT"), which are
// expected to be ISO-8601 / RFC 3339. It is the first option listed here:
//
// https://www.sqlite.org/datatype3.html#date_and_time_datatype:
//  - TEXT: "YYYY-MM-DD HH:MM:SS.SSS" (or with "T" in the blank position).
//  - REAL as Julian day numbers: the day count since 24 November 4714 BC.
//  - INTEGER as Unix time: the seconds count since 1970-01-01 00:00:00 UTC.
//
type TableConfig struct {
	TableName string
	ForenKeys []string
	Columns   []DbColSpec
}

// AllTableConfigs configures the three key tables.
var AllTableConfigs = []TableConfig{
	TableConfig_Inbatch,
	TableConfig_Contentity,
	TableConfig_Topicref,
}
