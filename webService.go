package main

import (
	"log"
	"net/http"
	"net/url"
)

type WebSvc struct {
}

func TransactionViewHandler(w http.ResponseWriter, r *http.Request, s *SharedExtConn) {

	log.Printf("INFO: Request From [%s] with Parameters [%s]\n", r.Header.Get("X-Real-IP"), r.URL)

	enableCors(&w)

	//query := make(url.Values)
	var err error
	_, err = url.ParseQuery(r.URL.RawQuery) // above 2 lines can be skipped by using ':='

	if err != nil {
		log.Println("ERROR: Request Error ", r.RemoteAddr, err, r.URL.RawQuery)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		log.Printf("ERROR: Processing Failed, Returning (%s)", "http.StatusBadRequest")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//check this
	w.Header().Set("Content-Type", "application/json")

	// write back to the client
	//_, err = fmt.Fprintf(w, "%s", status)
	//log.Printf("Status [%s]", status)

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

	log.Printf("Server Started [%s]", SERVER_PORT)
	http.ListenAndServe(string(SERVER_PORT), nil)
}
