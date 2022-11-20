package opensky

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	sevenDaysInHours  float64 = 24 * 7
	thirtyDaysInHours float64 = 24 * 30
	metersToFeet      float64 = 3.28084
	mPerSToFtPerMin   float64 = 196.85
	mPerSToKnots      float64 = 1.94384
)

// PositionSource determines whether source is ADS-B, ASTERIX, MLAT or FLARM
type PositionSource uint8

const (
	ADSB PositionSource = iota
	ASTERIX
	MLAT
	FLARM
)

func (p PositionSource) String() string {
	if p == 0 {
		return "ADS-B"
	}
	if p == 1 {
		return "ASTERIX"
	}
	if p == 2 {
		return "MLAT"
	}
	if p == 3 {
		return "FLARM"
	}
	return "Not a valid Position Source"
}

// Category represents the ICAO Vehicle category
type Category uint8

func (c Category) String() string {
	if c == 0 {
		return "N/A"
	}
	if c == 1 {
		return "No ADS-B Emitter Category Information"
	}
	if c == 2 {
		return "Light (< 15500 lbs)"
	}
	if c == 3 {
		return "Small (15500 to 75000 lbs)"
	}
	if c == 4 {
		return "Large (75000 to 300000 lbs)"
	}
	if c == 5 {
		return "High Vortex Large (aircraft such as B-757)"
	}
	if c == 6 {
		return "Heavy (> 300000 lbs)"
	}
	if c == 7 {
		return "High Performance (> 5g acceleration and 400 kts)"
	}
	if c == 8 {
		return "Rotorcraft"
	}
	if c == 9 {
		return "Glider / sailplane"
	}
	if c == 10 {
		return "Lighter-than-air"
	}
	if c == 11 {
		return "Parachutist / Skydiver"
	}
	if c == 12 {
		return "Ultralight / hang-glider / paraglider"
	}
	if c == 13 {
		return "Reserved"
	}
	if c == 14 {
		return "Unmanned Aerial Vehicle"
	}
	if c == 15 {
		return "Space / Trans-atmospheric vehicle"
	}
	if c == 16 {
		return "Surface Vehicle - Emergency Vehicle"
	}
	if c == 17 {
		return "Surface Vehicle - Service Vehicle"
	}
	if c == 18 {
		return "Point Obstacle (includes tethered balloons)"
	}
	if c == 19 {
		return "Cluster Obstacle"
	}
	if c == 20 {
		return "Line Obstacle"
	}
	return "Not a valid category"
}

type Sensors []int

type OskyValue interface {
	IsSet() bool
}

// Angle is a wrapepr type for degree values, such as track and heading
type Angle struct {
	float64
}

// Altitude is a wrapepr type for altitude related values, such as Geometric or Barometric Altitude
type Altitude struct {
	float64
}

// Speed is a wrapper type for speed related values, such as airspeed or vertical rate
type Speed struct {
	float64
}

func (a *Altitude) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		a.float64 = float64(value)
	default:
		a.float64 = -1
	}
	return nil
}

func (a *Angle) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		a.float64 = float64(value)
	default:
		a.float64 = -1
	}
	return nil
}

func (a *Speed) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		a.float64 = float64(value)
	default:
		a.float64 = -1
	}
	return nil
}

// StateVector contains the current state of an aircraft (or vehicle)
type StateVector struct {
	Icao24         string         `json:"icao24"`
	Callsign       string         `json:"callsign,omitempty"`
	OriginCountry  string         `json:"origin_country"`
	PositionTime   Timestamp      `json:"time_position,omitempty"`
	LastContact    Timestamp      `json:"last_contact,omitempty"`
	Longitude      Angle          `json:"longitude,omitempty"`
	Latitude       Angle          `json:"latitude,omitempty"`
	AltitudeBaro   Altitude       `json:"baro_altitude,omitempty"`
	OnGround       bool           `json:"on_ground"`
	Velocity       Speed          `json:"velocity,omitempty"`
	TrueTrack      Angle          `json:"true_track,omitempty"`
	VerticalRate   Speed          `json:"vertical_rate,omitempty"`
	Sensors        Sensors        `json:"sensors"`
	AltitudeGeo    Altitude       `json:"geo_altitude,omitempty"`
	Squawk         string         `json:"squawk,omitempty"`
	Spi            bool           `json:"spi"`
	PositionSource PositionSource `json:"position_source"`
	Category       Category       `json:"category"`
}

// StateVectorResponse contains a time and a set of StateVectors
type StateVectorResponse struct {
	Time   int64         `json:"time"`
	States []StateVector `json:"states"`
}

func degToString(val float64) string {
	return fmt.Sprintf("%.5f", val)
}

type FlightsResponse []Flight

type Timestamp struct {
	time.Time
}

func (t *Timestamp) Difference(other *Timestamp) time.Duration {
	return t.Sub(other.Time).Abs()
}

// Flight contains flight related values, such as airports
type Flight struct {
	Icao24                                string    `json:"icao24"`
	FirstSeen                             Timestamp `json:"firstSeen"`
	EstDepartureAirport                   string    `json:"estDepartureAirport"`
	LastSeen                              Timestamp `json:"lastSeen"`
	EstArrivalAirport                     string    `json:"estArrivalAirport"`
	Callsign                              string    `json:"callsign"`
	EstDepartureAirportHorizontalDistance int       `json:"estDepartureAirportHorizDistance"`
	EstDepartureAirportVerticalDistance   int       `json:"estDepartureAirportVertDistance"`
	EstArrivalAirportHorizontalDistance   int       `json:"estArrivalAirportHorizDistance"`
	EstArrivalAirportVerticalDistance     int       `json:"estArrivalAirportVertDistance"`
	DepartureAirportCandidates            int       `json:"departureAirportCandidatesCount"`
	ArrivalAirportCandidates              int       `json:"arrivalAirportCandidatesCount"`
}

// UnmarshalJSON unmarshal the raw json data into a given StateVector. Issue here is that provided JSON can contain something like
// "squawk": null, or "altitude": null, which is not handled well in go. Therefore we do manual checks here
func (s *StateVector) UnmarshalJSON(data []byte) error {

	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	s.Icao24, _ = v[0].(string)
	if v[1] != nil {
		s.Callsign, _ = v[1].(string)
	}
	s.OriginCountry = v[2].(string)

	if v[3] != nil {
		s.PositionTime = Timestamp{time.Unix(int64(v[3].(float64)), 0)}
	} else {
		s.PositionTime = Timestamp{time.Time{}}

	}

	s.LastContact = Timestamp{time.Unix(int64(v[4].(float64)), 0)}

	if v[5] != nil {
		val, _ := v[5].(float64)
		s.Longitude = Angle{val}
	}
	if v[6] != nil {
		val, _ := v[6].(float64)
		s.Latitude = Angle{val}
	}

	if v[7] != nil {
		val, _ := v[6].(float64)
		s.AltitudeBaro = Altitude{val}
	}
	s.OnGround = v[8].(bool)

	if v[9] != nil {
		val, _ := v[9].(float64)
		s.Velocity = Speed{val}
	}
	if v[10] != nil {
		val, _ := v[10].(float64)
		s.TrueTrack = Angle{val}
	}
	if v[11] != nil {
		val, _ := v[11].(float64)
		s.VerticalRate = Speed{val}
	}
	if v[12] != nil {
		s.Sensors = v[12].(Sensors)
	}
	if v[13] != nil {
		val, _ := v[13].(float64)
		s.AltitudeGeo = Altitude{val}
	}
	if v[14] != nil {
		s.Squawk = v[14].(string)
	}
	s.Spi = v[15].(bool)
	s.PositionSource = PositionSource(v[16].(float64))
	s.Category = Category(v[17].(float64))

	return nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		t.Time = time.Unix(int64(value), 0)
	default:
		t.Time = time.Time{} //zero time
	}
	return nil
}
