package main

import (
	"encoding/json"
	"log"
)

type JsonResponseData struct {
	Status  string      `json:"status"`
	Data    []MDataJSON `json:"data"`
	Errcode int         `json:"errcode"`
	Version string      `json:"version"`
}

func ResponseJsonDevice(status string, data []MDataJSON, errcode int) (result []byte, err error) {
	if TRACE {
		log.Printf("TRACE: Building Response JSON (status: [%s], data: [%+v], errcode: [%d]", status, data, errcode)
	}

	json_hash := &JsonResponseData{Status: status, Data: data, Errcode: errcode, Version: VERSION}

	result, err = json.Marshal(json_hash)
	if err != nil {
		log.Printf("ERROR: JSON Marshalling Failed [%s] data [%+v]", err, data)
		return
	}

	if TRACE {
		log.Printf("TRACE: Response JSON [%s]", string(result))
	}

	return
}
