package trams_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayubmalik/trams"
	"github.com/stretchr/testify/assert"
)

func TestClient_List(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Write([]byte(`{"Value": [{"Id": 1, "TLAREF": "BCH"}]}`))
			return
		}
		t.Errorf("unexpected method call %v", r.URL)
	}))

	t.Cleanup(func() {
		testServer.Close()
	})

	client := trams.NewClient(testServer.URL, 1000)

	t.Run("list all stations", func(t *testing.T) {
		metrolinks, err := client.List()

		if err != nil {
			assert.Fail(t, "got error", err)
		}
		assert.NotEmpty(t, metrolinks)
	})

}
