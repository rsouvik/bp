package main

type Transaction struct {
	cid   string
	image string
	descr string
	name  string
	S     *SharedExtConn // shared Connections
}

func NewTransaction(cid string, image string, descr string, name string, s *SharedExtConn) *Transaction {
	return &Transaction{cid: cid, image: image, descr: descr, name: name, S: s}
}

type TransactionJSON struct {
	image       string `json:"image"`
	description string `json:"description"`
	name        string `json:"name"`
}

type MData struct {
	cid   string
	image string
	descr string
	name  string
}
