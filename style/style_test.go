package style

import (
	"testing"
	"time"

	"github.com/ayubmalik/trams"
	"github.com/stretchr/testify/assert"
)

func TestStyle_FormatMetrolink(t *testing.T) {

	m := trams.Metrolink{
		Id:              699,
		Line:            "Airport",
		TLAREF:          "BCH",
		PIDREF:          "BCH-TPID01",
		StationLocation: "Benchill",
		AtcoCode:        "9400ZZMABLL1",
		Direction:       "Outgoing",
		Dest0:           "Manchester Airport",
		Carriages0:      "Single",
		Status0:         "Due",
		Wait0:           "12",
		Dest1:           "Manchester Airport",
		Carriages1:      "Single",
		Status1:         "Due",
		Wait1:           "18",
		Dest2:           "",
		Carriages2:      "",
		Status2:         "",
		Wait2:           "",
		Dest3:           "",
		Carriages3:      "",
		Status3:         "",
		MessageBoard:    "Services are now able to run through Piccadilly. Thank you for your patience during this time and we apologise for any inconvenience this may have caused to your journey this morning",
		Wait3:           "",
		LastUpdated:     time.Now(),
	}

	formatted := FormatMetrolink(m, 1)
	assert.NotEmpty(t, formatted)
}
