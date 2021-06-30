package trams

import (
	"encoding/json"
	"net/http"
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

// List retrieves information about the specified Metrolink stations.
// If no station IDs is empty retrives all stations.
func (c Client) List(ids ...string) ([]Metrolink, error) {
	query := ""
	if len(ids) > 0 {
		query += "?id=" + strings.Join([]string(ids), "&id=")
	}

	resp, err := http.Get(c.url + query)
	if err != nil {
		return nil, err
	}

	var metrolinks []Metrolink
	err = json.NewDecoder(resp.Body).Decode(&metrolinks)
	if err != nil {
		return nil, err
	}

	return metrolinks, err
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

/*
sample API response
{
    "Id": 699,
    "Line": "Airport",
    "TLAREF": "BCH",
    "PIDREF": "BCH-TPID01",
    "StationLocation": "Benchill",
    "AtcoCode": "9400ZZMABLL1",
    "Direction": "Outgoing",
    "Dest0": "Manchester Airport",
    "Carriages0": "Single",
    "Status0": "Due",
    "Wait0": "12",
    "Dest1": "Manchester Airport",
    "Carriages1": "Single",
    "Status1": "Due",
    "Wait1": "18",
    "Dest2": "",
    "Carriages2": "",
    "Status2": "",
    "Wait2": "",
    "Dest3": "",
    "Carriages3": "",
    "Status3": "",
    "MessageBoard": "Services are now able to run through Piccadilly. Thank you for your patience during this time and we apologise for any inconvenience this may have caused to your journey this morning",
    "Wait3": "",
    "LastUpdated": "2021-06-24T13:26:38Z"
  }
*/
