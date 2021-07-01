package trams_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayubmalik/trams"
	"github.com/stretchr/testify/assert"
)

func TestClient_Get(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := r.URL.Query().Get("id"); id != "" {
			w.Write([]byte(`[{"Id": 1}]`))
			return
		}

		if r.URL.Path == "/" {
			w.Write([]byte(`[{"Id": 1}, {"Id": 2}]`))
			return
		}

		t.Errorf("unexpected method call %v", r.URL)
	}))

	t.Cleanup(func() {
		testServer.Close()
	})

	client := trams.NewClient(testServer.URL, 1000)

	t.Run("get all stations", func(t *testing.T) {
		metrolinks, err := client.Get()

		if err != nil {
			assert.Fail(t, "got error", err)
		}
		assert.Len(t, metrolinks, 2)
	})

	t.Run("get specific by IDs", func(t *testing.T) {
		metrolinks, err := client.Get("1", "2")

		if err != nil {
			assert.Fail(t, "got error", err)
		}
		assert.Len(t, metrolinks, 1)
	})

}
