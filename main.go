package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var (
	p = NewPinger("chaostreff.vpn.zekjur.net", time.Minute)
)

func HandleGet(res http.ResponseWriter, req *http.Request) {
	ep := NewEndpoint()
	ep.State.Open = p.GetState()

	enc := json.NewEncoder(res)
	err := enc.Encode(ep)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/spaceapi.json", HandleGet)

	err := http.ListenAndServe("localhost:5124", nil)
	if err != nil {
		log.Fatal(err)
	}
}
