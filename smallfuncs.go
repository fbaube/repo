package repo

// Times has (create, import, last edit) and uses only ISO-8601 / RFC 3339.
type Times struct {
	T_Cre string
	T_Imp string
	T_Edt string
}

func checkerr(e error) {
	if e == nil {
		return
	}
	panic("Sqlite3 FAILURE: " + e.Error())
}

func (p *SimpleRepo) MustExecStmt(s string) {
	stmt, e := p.Handle().Prepare(s)
	checkerr(e)
	_, e = stmt.Exec() // rslt,e := ...
	checkerr(e)
	// liid, _ := rslt.LastInsertId()
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("DD:mustExecStmt: ID %d nR %d \n", liid, naff)
}
