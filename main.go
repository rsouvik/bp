package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {

	sec := &SharedExtConn{}

	// make mysql connection
	// build the connection string user:passwd@tcp(host:port)/db
	//sql_conn_str := fmt.Sprintf("postgres://%s:%s@tcp(%s:%d)/%s", DB_USER, DB_PASSWD, DB_HOST, DB_PORT, DB_DATABASE)
	sql_conn_str := fmt.Sprintf("postgres://%s:%s@%s/%s", DB_USER, DB_PASSWD, DB_HOST, DB_DATABASE)
	sql_ctx, err := GetSqlConn(sql_conn_str)
	if err != nil {
		log.Fatalf("FATAL: PG Connection Failed! [%s]", err)
	}

	// save the mysql connection context
	sec.Msql = sql_ctx

	//test
	//transaction := NewTransaction("abda", "abca", "abca", "abca" /*,sec*/)
	//_, err = sec.Msql.InsertTransaction(transaction)
	//if err != nil {
	//	log.Fatalf("FATAL: PG Insert Failed! [%s]", err)
	//}

	doneCh := make(chan struct{})
	//mdataChannel := make(chan MData, 2)

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	log.Println("Starting Service!")

	//test
	//cids := []string{"bafybeia67q6eabx2rzu6datbh3rnsoj7cpupudckijgc5vtxf46zpnk2t4/3885"}

	var cids []string

	filePath := os.Args[1]
	readFile, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		cids = append(cids, fileScanner.Text())
	}

	readFile.Close()

	//Start service
	scrapeSvc := NewScrapeSvc()
	go scrapeSvc.Run(cids, doneCh, sec)

	//Webservice
	webSvc := WebSvc{}
	go webSvc.Run(sec)

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
