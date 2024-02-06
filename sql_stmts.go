package main

import (
	_ "database/sql"
	_ "github.com/lib/pq"
	"log"
)

//TBD
func (sql_ctx *SqlContext) InsertTransaction(tr *Transaction) (id int64, err error) {
	//if DEBUG {
	//	log.Printf("DEBUG: InsertTransaction Hash %v", hash)
	//}

	stmtIns, err := sql_ctx.Db.Prepare("INSERT INTO cidmeta (cid,image,descr,ciname) VALUES($1,$2,$3,$4)")
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	// execute insert
	res, err := stmtIns.Exec(tr.cid, tr.image, tr.descr, tr.name)
	if err != nil {
		log.Printf("ERROR: Insert INTO transactions (%s)", err)
		return -1, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		log.Printf("ERROR: Retrieving LastInsertId Failed (%s)", err)
		return -1, err
	}

	return
}
