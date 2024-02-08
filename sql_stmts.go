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
	_, err = stmtIns.Exec(tr.cid, tr.image, tr.descr, tr.name)
	if err != nil {
		log.Printf("ERROR: Insert INTO transactions (%s)", err)
		return -1, err
	}

	/*id, err = res.LastInsertId()
	if err != nil {
		log.Printf("ERROR: Retrieving LastInsertId Failed (%s)", err)
		return -1, err
	}*/

	return -2, nil
}

func (sql_ctx *SqlContext) GetMetaData(cid string) ([]*MData, error) {

	rows, err := sql_ctx.Db.Query("SELECT cid,image,descr,ciname FROM cidmeta where cid = $1", cid)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	mdatas := make([]*MData, 0)

	for rows.Next() {
		md := new(MData)
		if err := rows.Scan(&md.cid, &md.image, &md.descr, &md.name); err != nil {
			panic(err)
		}
		mdatas = append(mdatas, md)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return mdatas, nil

}

func (sql_ctx *SqlContext) GetMetaDataAll() ([]*MData, error) {

	rows, err := sql_ctx.Db.Query("SELECT cid,image,descr,ciname FROM cidmeta")
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	mdatas := make([]*MData, 0)

	for rows.Next() {

		md := new(MData)
		if err := rows.Scan(&md.cid, &md.image, &md.descr, &md.name); err != nil {
			panic(err)
		}
		//log.Println(md.name)
		mdatas = append(mdatas, md)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return mdatas, nil

}
