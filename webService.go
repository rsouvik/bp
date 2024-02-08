package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type WebSvc struct {
}

func TransactionViewHandlerCID(w http.ResponseWriter, r *http.Request, s *SharedExtConn) {

	log.Printf("INFO: Request From [%s] with Parameters [%s]\n", r.Header.Get("X-Real-IP"), r.URL)

	enableCors(&w)

	query := make(url.Values)
	var err error
	query, err = url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		log.Println("ERROR: Request Error ", r.RemoteAddr, err, r.URL.RawQuery)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var mdatas []*MData

	// map to store flattened struct
	var q = make(map[string]string)
	// flatten the GET params
	for k, v := range query {
		q[k] = v[0] // we might have more than 1, but we stick to first one
	}

	log.Printf("q[cid] is (%s)", q["cid"])

	mdatas, _ = s.Msql.GetMetaData(q["cid"])

	var tr []MDataJSON
	for k := 0; k < len(mdatas); k++ {
		tr = append(tr, MDataJSON{mdatas[k].cid, mdatas[k].image, mdatas[k].descr, mdatas[k].name})
	}

	msg, err := ResponseJsonDevice("SUCCEEDED", tr, 0)

	status := ""
	if err != nil {
		status = err.Error()
	}
	//else {
	//	status = string(msg)
	//}

	if err != nil {
		log.Printf("ERROR: Processing Failed, Returning (%s)", "http.StatusBadRequest")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//check this
	w.Header().Set("Content-Type", "application/json")

	prettyJSON, err := formatJSON(msg)
	if err != nil {
		log.Fatal(err)
	}

	status = string(prettyJSON)

	// write back to the client
	_, err = fmt.Fprintf(w, "%s", status)
	log.Printf("Status [%s]", status)

	return

}

func TransactionViewHandler(w http.ResponseWriter, r *http.Request, s *SharedExtConn) {

	log.Printf("INFO: Request From [%s] with Parameters [%s]\n", r.Header.Get("X-Real-IP"), r.URL)

	enableCors(&w)

	var mdatas []*MData
	mdatas, _ = s.Msql.GetMetaDataAll()

	var tr []MDataJSON
	for k := 0; k < len(mdatas); k++ {
		tr = append(tr, MDataJSON{mdatas[k].cid, mdatas[k].image, mdatas[k].descr, mdatas[k].name})
	}

	msg, err := ResponseJsonDevice("SUCCEEDED", tr, 0)

	status := ""
	if err != nil {
		status = err.Error()
	}
	//else {
	//	status = string(msg)
	//}

	if err != nil {
		log.Printf("ERROR: Processing Failed, Returning (%s)", "http.StatusBadRequest")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//check this
	w.Header().Set("Content-Type", "application/json")

	prettyJSON, err := formatJSON(msg)
	if err != nil {
		log.Fatal(err)
	}

	status = string(prettyJSON)
	// write back to the client
	_, err = fmt.Fprintf(w, "%s", status)
	log.Printf("Status [%s]", status)

	return

}

// This is a generic request Handler. It returns a function
// which creates a "go" routine for each kind of request
func genericHandler(fn func(http.ResponseWriter, *http.Request, *SharedExtConn), sec *SharedExtConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { // returns void, we don't care!
		fn(w, r, sec)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (p *WebSvc) Run(s *SharedExtConn) {

	// start the webserver
	http.HandleFunc("/token", genericHandler(TransactionViewHandler, s))

	http.HandleFunc("/token/cid", genericHandler(TransactionViewHandlerCID, s))

	log.Printf("Server Started [%s]", SERVER_PORT)
	http.ListenAndServe(string(SERVER_PORT), nil)
}
