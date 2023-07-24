package trams

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

// NewClient returns a new trams client where url is the backend cloud
// function URL and timeout is the timeout in milliseconds for responses.
func NewClient(url string, timeout int) Client {
	return Client{url: url, timeout: timeout}
}

// Client is the main API for communicating with the trams backend cloud function.
// TODO: use timeout
type Client struct {
	url     string
	timeout int
}

// Get retrieves information about the specified Metrolink stations.
// If no station IDs is empty retrives all stations.
func (c Client) Get(ids ...string) ([]Metrolink, error) {
	query := ""
	if len(ids) > 0 {
		query += "?id=" + strings.Join([]string(ids), "&id=")
	}

	resp, err := http.Get(c.url + query)
	if err != nil {
		return nil, err
	}

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("TRAMS_DEBUG") == "1" {
		fmt.Println("--- Start Response Body (Get)")
		fmt.Printf("%q\n", dump)
		fmt.Println("--- End Response Body (Get)")
	}

	var metrolinks []Metrolink
	err = json.NewDecoder(resp.Body).Decode(&metrolinks)
	if err != nil {
		return nil, err
	}

	return metrolinks, err
}

// List all available station IDs.
func (c Client) List() ([]StationID, error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return nil, err
	}

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("TRAMS_DEBUG") == "1" {
		fmt.Println("--- Start Response Body (List)")
		fmt.Printf("%q\n", dump)
		fmt.Println("--- End Response Body (List)")
	}

	var stationIDs []StationID
	err = json.NewDecoder(resp.Body).Decode(&stationIDs)
	if err != nil {
		return nil, err
	}

	return stationIDs, err
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

func (m Metrolink) Platform() string {
	return m.PIDREF[len(m.PIDREF)-1:]
}

// StationID identifies a Metrolink station.
type StationID struct {
	Id              int
	TLAREF          string
	StationLocation string
}
