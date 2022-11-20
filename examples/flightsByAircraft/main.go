package main

import (
	"log"
	"time"

	"github.com/chrisdo/opensky-go-api"
)

func main() {

	client := opensky.NewClient()

	begin, end := opensky.GetStartAndEndOfDay(time.Now().UTC())
	flightsReq, err := opensky.NewFlightsByAircraftRequest("3c66e5", begin, end)
	if err != nil {
		log.Println(err)
	}
	flights, err := client.RequestFlightsByAircraft(flightsReq)
	if err != nil {
		log.Println(err)
	}
	for _, f := range *flights {
		log.Printf("%v\n", f)
	}
}
