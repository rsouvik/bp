package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

// later we can use to store more context data
type SqlContext struct {
	Db *sql.DB
}

func GetSqlConn(connStr string) (*SqlContext, error) {
	db, err := sql.Open("postgres", connStr) // user:password@/database OR user:passwd@tcp(host:port)/db
	if err != nil {
		db.Close()
		log.Println(err.Error())
		return nil, err
	}

	sql_ctx := &SqlContext{}
	sql_ctx.Db = db

	return sql_ctx, err
}

func (sql_ctx *SqlContext) CloseSqlConn() (err error) {
	err = sql_ctx.Db.Close()
	return
}
