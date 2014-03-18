package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var (
	pinger    = NewPinger("chaostreff.vpn.zekjur.net", time.Minute)
	locpoller = NewLocationPoller(10 * time.Minute)
)

func HandleGet(res http.ResponseWriter, req *http.Request) {
	ep := NewEndpoint()
	ep.Location = locpoller.Get()
	if ep.Location == uniLocation {
		ep.State.Open = pinger.GetState()
	} else {
		now := time.Now()
		if now.Weekday() == time.Thursday && now.Day() < 8 && now.Hour() >= 19 && now.Hour() < 22 {
			ep.State.Open = True
		} else {
			ep.State.Open = False
		}
	}

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
