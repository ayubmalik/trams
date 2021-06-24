package trams

import (
	"time"
)

// Client is the main API for communicating with the trams backend cloud function.
type Client struct {

	// URL of the backend cloud function.
	URL string

	// Timeout in milliseconds for receiving result from backend
	Timeout int
}

// Metrolinks is the JSON struct returned by the backend cloud function.
type Metrolinks struct {
	Value []Metrolink
}

// Metrolink provides information for a station location.
type Metrolink struct {
	Id              int
	Line            string
	TLAREF          string
	PIDREF          string
	StationLocation string
	AtcoCode        string
	Direction       string
	Dest0           string
	Carriages0      string
	Status0         string
	Wait0           string
	Dest1           string
	Carriages1      string
	Status1         string
	Wait1           string
	Dest2           string
	Carriages2      string
	Status2         string
	Wait2           string
	Dest3           string
	Carriages3      string
	Status3         string
	MessageBoard    string
	Wait3           string
	LastUpdated     time.Time
}
