package opensky

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	baseUrl                      = "https://opensky-network.org/api/"
	allStates             string = baseUrl + "states/all"
	ownStates             string = baseUrl + "states/own"
	flightsWithinInterval string = baseUrl + "flights/all"
	FlightsByAircraft     string = baseUrl + "flights/aircraft"
	arrivalsByAirport     string = baseUrl + "flights/arrival"
	departuresByAirport   string = baseUrl + "flights/departure"
)

type BoundingBox struct {
	latMin float64
	latMax float64
	lonMin float64
	lonMax float64
}

func NewBoundingBox(latMin, latMax, lonMin, lonMax float64) *BoundingBox {
	return &BoundingBox{latMin, latMax, lonMin, lonMax}
}

type AirportRequestType int

const (
	Departure AirportRequestType = iota
	Arrival
)

// StateVectorRequest represents a query to retrieva all state vectors for given parameters, see also
// https://openskynetwork.github.io/opensky-api/rest.html#all-state-vectors
type StateVectorRequest struct {
	url.Values
}

func NewStateVectorRequest() *StateVectorRequest {
	return &StateVectorRequest{make(url.Values)}
}

func (r *StateVectorRequest) WithBoundingBox(b *BoundingBox) *StateVectorRequest {
	r.Add("lamin", degToString(b.latMin))
	r.Add("lomin", degToString(b.lonMin))
	r.Add("lamax", degToString(b.latMax))
	r.Add("lomax", degToString(b.lonMax))
	return r
}

func (r *StateVectorRequest) WithIcaoAdress(icao24 string) *StateVectorRequest {
	if len(icao24) != 6 {
		log.Println("Length of icao must be 6")
	}
	r.Add("icao24", icao24)
	return r
}

func (r *StateVectorRequest) AtTime(time time.Time) *StateVectorRequest {
	r.Add("time", fmt.Sprint(time.Unix()))
	return r
}

func (r *StateVectorRequest) IncludeCategory() *StateVectorRequest {
	r.Add("extended", strconv.Itoa(1))
	return r
}

// OwnStateVectorsRequest can be used to retrieve all state vectors for your own sensors, or filtered by parameters
// This request only works when username and password is provided
// See also: https://openskynetwork.github.io/opensky-api/rest.html#own-state-vectors
type OwnStateVectorsRequest struct {
	url.Values
}

func (r *OwnStateVectorsRequest) WithSensors(sensors ...int) *OwnStateVectorsRequest {
	for _, sensor := range sensors {
		r.Add("serials", fmt.Sprint(sensor))
	}
	return r
}

func (r *OwnStateVectorsRequest) WithIcaoAdress(icao24 string) *OwnStateVectorsRequest {
	if len(icao24) != 6 {
		log.Println("Length of icao must be 6")
	}
	r.Add("icao24", icao24)

	return r
}

func (r *OwnStateVectorsRequest) AtTime(time time.Time) *OwnStateVectorsRequest {
	r.Add("time", fmt.Sprint(time.Unix()))
	return r
}

// FlightsWithinIntervalRequest can be used to retrieved Flights within a given interval
// The interval must not be larger than 2 hours
// https://openskynetwork.github.io/opensky-api/rest.html#flights-in-time-interval
type FlightsWithinIntervalRequest struct {
	url.Values
}

func NewFlightsWithinIntervallRequest(begin, end time.Time) (*FlightsWithinIntervalRequest, error) {
	err := checkInterval(begin, end, 2)
	if err != nil {
		return nil, err
	}

	r := FlightsWithinIntervalRequest{make(url.Values)}
	r.Add("begin", fmt.Sprint(begin.Unix()))
	r.Add("end", fmt.Sprint(end.Unix()))

	return &r, nil
}

// FlightsByAirportRequest can be used to retrieve all flights for a given airport interval. The interval must not exceed 7 days
// whether departures or arrivals are of interest is determined by the provided paramter AirportReuqestType
// https://openskynetwork.github.io/opensky-api/rest.html#arrivals-by-airport
type FlightsByAirportRequest struct {
	url.Values
	t AirportRequestType
}

func NewFlightsByAirportRequest(airport string, begin time.Time, end time.Time, t AirportRequestType) (*FlightsByAirportRequest, error) {
	if len(airport) != 4 {
		return nil, errors.New("ICAO airport must be 4 characters long")
	}
	err := checkInterval(begin, end, sevenDaysInHours)
	if err != nil {
		return nil, err
	}
	r := FlightsByAirportRequest{make(url.Values), t}
	r.Add("airport", strings.ToUpper(airport))
	r.Add("begin", fmt.Sprint(begin.Unix()))
	r.Add("end", fmt.Sprint(end.Unix()))

	return &r, nil
}

// FlightsByAircraftRequest can be used to get flights for a given aircraft within a given interval.
// the interval must not exceed 30 days
// https://openskynetwork.github.io/opensky-api/rest.html#flights-by-aircraft
type FlightsByAircraftRequest struct {
	url.Values
}

func NewFlightsByAircraftRequest(icao24 string, begin time.Time, end time.Time) (*FlightsByAircraftRequest, error) {
	if len(icao24) != 6 {
		return nil, errors.New("icao24 address must be exactly 6 characters")
	}
	err := checkInterval(begin, end, thirtyDaysInHours)
	if err != nil {
		return nil, err
	}

	r := FlightsByAircraftRequest{make(url.Values)}
	r.Add("icao24", icao24)
	r.Add("begin", fmt.Sprint(begin.Unix()))
	r.Add("end", fmt.Sprint(end.Unix()))

	return &r, nil
}

func checkInterval(begin, end time.Time, limitInHours float64) error {
	if end.Before(begin) {
		return errors.New("end must be AFTER begin")
	}
	if end.Sub(begin).Abs().Hours() > limitInHours {
		return errors.New("interval duration must not exceed 7 days")
	}
	return nil
}
