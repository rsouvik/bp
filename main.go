package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {

	sec := &SharedExtConn{}

	// make mysql connection
	// build the connection string user:passwd@tcp(host:port)/db
	log.Println(DB_USER)
	sql_conn_str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", DB_USER, DB_PASSWD, DB_HOST, DB_PORT, DB_DATABASE)
	sql_ctx, err := GetSqlConn(sql_conn_str)
	if err != nil {
		log.Fatalf("FATAL: PG Connection Failed! [%s]", err)
	}

	// save the mysql connection context
	sec.Msql = sql_ctx
	log.Println("INFO: PGSql Connection Context - OK")

	//test
	transaction := NewTransaction("abd", "abc", "abc", "abc", sec)
	_, err = sec.Msql.InsertTransaction(transaction)
	if err != nil {
		log.Fatalf("FATAL: PG Insert Failed! [%s]", err)
	}

	doneCh := make(chan struct{})

	//should this be here ?
	defer sec.Msql.CloseSqlConn()

	go Ctrl(doneCh)

	<-doneCh
	log.Println("Shutting down!")

}

// Ctrl handles monitor shutdown actions.
func Ctrl(doneChan chan<- struct{}) {
	//log.Println("Shutting down!")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	log.Println("Shutting down interrupt!")
	for range sigChan {
		doneChan <- struct{}{}
	}
}
