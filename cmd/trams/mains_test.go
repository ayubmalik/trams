package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/ayubmalik/trams"
	"github.com/stretchr/testify/assert"
)

func TestReadCachedStations(t *testing.T) {

	cache := path.Join(os.TempDir(), "test-trams-cache")

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `[{"TLAREF":"REF1","StationLocation":"Station From Client"}]`)
	}))

	t.Cleanup(func() {
		testServer.Close()
		os.RemoveAll(cache)
	})

	client := trams.NewClient(testServer.URL, 1000)

	t.Run("calls Client.List() when cache does not exist", func(t *testing.T) {
		stations, err := cachedStations(client, cache)
		if err != nil {
			assert.FailNow(t, "error", err)
		}

		assert.Len(t, stations, 1)
		assert.Equal(t, "Station From Client", stations["REF1"][0].StationLocation)
		assert.FileExists(t, cache)
	})

	t.Run("does not call Client.List() when cache exists", func(t *testing.T) {
		f, err := os.Create(cache)
		if err != nil {
			assert.FailNow(t, "could not create cache")
		}
		io.WriteString(f, `[{"TLAREF":"REF1","StationLocation":"Station From Cache"}]`)

		t.Cleanup(func() {
			f.Close()
		})

		stations, err := cachedStations(client, cache)
		if err != nil {
			assert.FailNow(t, "error reading stations from cache", err)
		}

		// assert.Len(t, stations, 1)
		fmt.Println("YYYY", stations)
		assert.Equal(t, "Station From Cache", stations["REF1"][0].StationLocation)
	})
}
