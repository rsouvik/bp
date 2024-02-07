package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_ "strings"
	"time"
)

type ScrapeSvc struct {
	DbCtrl       *DbCtrl
	statusMap    map[string]bool
	mdataChannel chan string
}

func NewScrapeSvc() *ScrapeSvc {
	return &ScrapeSvc{DbCtrl: NewDbCtrl(),
		statusMap: make(map[string]bool), mdataChannel: make(chan string)}
}

func (p *ScrapeSvc) fetchMeta(done chan struct{}, s *SharedExtConn) {

processingLoop:

	for {
		select {

		case cid := <-p.mdataChannel:

			resp, err := http.Get("https://ipfs.io/ipfs/" + cid)
			if err != nil {
				log.Fatalln(err)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			var data map[string]interface{}
			err = json.Unmarshal([]byte(body), &data)
			if err != nil {
				panic(err)
			}

			tr := NewTransaction(cid, data["image"].(string), data["description"].(string), data["name"].(string))

			s.Msql.InsertTransaction(tr)

			time.Sleep(30 * time.Second)

		case <-done:

			break processingLoop

		}
	}

}

func (p *ScrapeSvc) Run(cids []string, done chan struct{}, s *SharedExtConn) error {

	log.Println("Scraping Service Started!")

	//Read through CID array and fetch
	for i := 0; i < len(cids); i++ {
		p.mdataChannel <- cids[i]
	}

	concurrentThreads := 3
	//Write to channel
	for i := 0; i < concurrentThreads; i++ {
		go p.fetchMeta(done, s)
	}

processingLoop:
	for {
		select {

		case <-done:

			break processingLoop

		}
	}

	return nil

}
