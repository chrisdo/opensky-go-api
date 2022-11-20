package main

import (
	"log"
	"time"

	"github.com/chrisdo/opensky-go-api"
)

func main() {

	client := opensky.NewClient()

	log.Println()

	flightsByirport(client, "EDDF", opensky.Arrival)

	flightsByirport(client, "EDDF", opensky.Departure)
}

func flightsByirport(client *opensky.Client, airport string, t opensky.AirportRequestType) {
	begin, end := opensky.GetStartAndEndOfDay(time.Now().UTC())
	airportsReq, err := opensky.NewFlightsByAirportRequest(airport, begin, end, t)
	if err != nil {
		log.Println(err)
	}
	flights, err := client.RequestFlightsByAirport(airportsReq)
	if err != nil {
		log.Println(err)
	}
	for _, f := range *flights {
		log.Printf("%v\n", f)
	}
}
