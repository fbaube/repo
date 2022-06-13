package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	FP "path/filepath"

	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
	XU "github.com/fbaube/xmlutils"
)

// NewContentityRecord works for directories and symlinks too.
// It used to SetError(..), but no longer does (not much anyways).
func NewContentityRecord(pPP *FU.PathProps) (*ContentityRecord, error) {
	var e error
	pCR := new(ContentityRecord)
	pCR.PathProps = *pPP

	if !pPP.Exists() {
		return pCR, errors.New("does not exist")
	}
	if pPP.Size() == 0 {
		// panic("barf")
		return pCR, nil
	}
	if pPP.IsOkayDir() || pPP.IsOkaySymlink() {
		// COMMENTING THIS OUT IS A FIX
		// pCR.SetError(errors.New("Is directory or symlink"))
		return pCR, nil
	}
	if !pPP.IsOkayFile() {
		return pCR, errors.New("is not valid file")
	}
	// OK, it's a valid file.
	// But maybe it's empty!
	if pPP.Size() == 0 {
		L.L.Progress("Skipping fetch for zero-length content")
	} else {
		pCR.ContentityStructure.Raw, e = pPP.FetchContent()
		if e != nil {
			L.L.Error("DB.newCnty: cannot fetch content: " + e.Error())
			return pCR, fmt.Errorf("DB.newCnty: cannot fetch content: %w", e)
		}
	}
	var pAR *XU.AnalysisRecord
	pAR, e = FU.AnalyseFile(pCR.ContentityStructure.Raw, FP.Ext(string(pPP.AbsFP)))
	if e != nil {
		L.L.Error("DB.newCnty: analyze file failed: " + e.Error())
		return pCR, fmt.Errorf("fu.newCR: analyze file failed: %w", e)
	}
	if pAR.MType == "" {
		L.L.Warning("No MType, so trying snift-MIME-type: %s", pAR.MimeTypeAsSnift)
		switch pAR.MimeTypeAsSnift {
		case "text/xml/image/svg+xml":
			println("SVG!!")
			pAR.MType = "xml/cnt/svg"
		}
	}
	if pAR == nil {
		panic("NIL pAR")
	}
	pCR.AnalysisRecord = *pAR
	// SPLIT FILE!
	if !pAR.ContentityStructure.HasNone() {
		L.L.Okay("Key elm triplet: Root<%s> Meta<%s> Text<%s>",
			pAR.ContentityStructure.Root.String(),
			pAR.ContentityStructure.Meta.String(),
			pAR.ContentityStructure.Text.String())
	} else if pAR.FileType() == "MKDN" {
		// pAR.KeyElms.SetToAllText()
		// L.L.Warning("TODO set MKDN all text, and ranges")
	} else if pAR.FileType() == "BIN" {
	} else {
		L.L.Warning("Found no key elms (root,meta,text)")
	}
	// fmt.Printf("D=> NewCR: %s \n", pCR.String())
	return pCR, nil
}

// GetContentityAll gets all content in the DB.
func (p *SimpleRepo) GetContentityAll() (pp []*ContentityRecord) {
	pp = make([]*ContentityRecord, 0, 16)
	rows, err := p.Handle().Queryx("SELECT * FROM CONTENT")
	if err != nil {
		panic("GetContentityAll")
	}
	for rows.Next() {
		p := new(ContentityRecord)
		err := rows.StructScan(p)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("    DD:%#v\n", *p)
		pp = append(pp, p)
	}
	return pp
}

// InsertContentityRecord adds a content item (i.e. a file) to the DB.
func (p *SimpleRepo) InsertContentityRecord(pC *ContentityRecord) (int, error) {
	var rslt sql.Result
	var stmt string
	var e error
	// println("REL:", pC.RelFP)
	// println("ABS:", pC.AbsFP)

	pC.T_Cre = SU.Now() // time.Now().UTC().Format(time.RFC3339)
	pC.T_Imp = SU.Now() // time.Now().UTC().Format(time.RFC3339)
	tx := p.Handle().MustBegin()
	stmt = "INSERT INTO CONTENTITY(" +
		"idx_inbatch, descr, relfp, absfp, " +
		"t_cre, t_imp, t_edt, " +
		// "metaraw, textraw, " +
		"mimetype, mtype, " +
		// roottag, rootatts, " +
		"xmlcontype, xmldoctype, ditaflavor, ditacontype" +
		") VALUES(" +

		// ":idx_inbatch, :pathprops.relfp, :pathprops.absfp, " +
		":idx_inbatch, :descr, :relfp, :absfp, " +

		// ":times.t_cre, :times.t_imp, :times.t_edt, " +
		":t_cre, :t_imp, :t_edt, " +

		// ":metaraw, :textraw, " +
		// ":mimetype, :mtype, " +
		":mimetype, :mtype, " +
		// ":root.name, :root.atts, " +
		// ":analysisrecord.contentitystructure.root.name, " +
		// ":analysisrecord.contentitystructure.root.atts, " +

		":xmlcontype, :doctype, :ditaflavor, :ditacontype);"
		// ":doctype, :ditaflavor, :ditacontype);"

	rslt, e = tx.NamedExec(stmt, pC)
	tx.Commit()
	// println("=== ### ===")
	if e != nil {
		L.L.Error("DB.Add_Contentity: %s", e.Error())
	}
	if e != nil {
		println("========")
		println("DB: NamedExec: ERROR:", e.Error())
		println("========")
		fnam := "./insert-Contentity-failed.sql"
		e = ioutil.WriteFile(fnam, []byte(stmt), 0644)
		if e != nil {
			L.L.Error("Could not write file: " + fnam)
		} else {
			L.L.Dbg("Wrote \"INSERT INTO contentity ... \" to: " + fnam)
		}
		panic("INSERT CONTENTITY failed")
	}
	liid, _ := rslt.LastInsertId()
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("    DD:InsertFile: ID %d nR %d \n", liid, naff)
	return int(liid), nil
}
