package opensky

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	http.Client
	username string
	password string
}

// NewClient creates a new opensky client with no username and password reslting in anonymous mode. Opensky restrictions apply
func NewClient() *Client {
	return &Client{
		Client:   http.Client{Timeout: time.Second * 10},
		username: "",
		password: "",
	}
}

// NewRegisteredClient creates an ew client with povided username and password
func NewRegisteredClient(username, password string) *Client {
	return &Client{
		Client:   http.Client{Timeout: time.Second * 10},
		username: username,
		password: password,
	}
}

func (c *Client) RequestStateVectors(r *StateVectorRequest) (*StateVectorResponse, error) {
	var result StateVectorResponse
	err := doRequest(c, allStates, r.Values, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) RequestOwnStateVectors(r *OwnStateVectorsRequest) (*StateVectorResponse, error) {
	var result StateVectorResponse
	err := doRequest(c, ownStates, r.Values, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) RequestFlightsByAircraft(r *FlightsByAircraftRequest) (*FlightsResponse, error) {

	var result FlightsResponse
	err := doRequest(c, FlightsByAircraft, r.Values, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) RequestFlightsWithinInterval(r *FlightsWithinIntervalRequest) (*FlightsResponse, error) {

	var result FlightsResponse
	err := doRequest(c, flightsWithinInterval, r.Values, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) RequestFlightsByAirport(r *FlightsByAirportRequest) (*FlightsResponse, error) {

	var url string
	if r.t == Departure {
		url = departuresByAirport
	} else if r.t == Arrival {
		url = arrivalsByAirport
	} else {
		return nil, errors.New("must provide either Departure or Arrival request type")
	}
	var result FlightsResponse
	err := doRequest(c, url, r.Values, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil

}

func doRequest(c *Client, url string, params url.Values, response interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	req.URL.RawQuery = params.Encode()
	log.Printf("Requesting: %s\n", req.URL.String())

	res, err := c.Do(req)
	if err != nil {
		return err
	}
	log.Printf("Response received with status: %s\n", res.Status)
	if res.StatusCode == 403 {
		return errors.New("not allowed to make this request without username and password")
	}

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &response)
	if err != nil {
		return err
	}
	return nil
}
