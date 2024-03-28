package sync

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExchangeRateSync_loadHTTPData(t *testing.T) {
	// Create a test server to mock the HTTP response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a sample XML response
		xmlResponse := `<ExchangeRate>
			<Currency>USD</Currency>
			<Rate>1.23</Rate>
		</ExchangeRate>`
		w.Write([]byte(xmlResponse))
	}))
	defer server.Close()

	// Create an instance of ExchangeRateSync with the test server URL
	sync := &ExchangeRateSync{
		url:        server.URL,
		httpClient: server.Client(),
	}

	t.Run("valid HTTP response", func(t *testing.T) {
		// Call the loadHTTPData method
		data, err := sync.loadHTTPData()

		// Assert that no error occurred
		assert.NoError(t, err)

		// Assert that the data is correctly parsed
		expectedData := ExchangeRate{
			Currency: "USD",
			Rate:     1.23,
		}
		assert.Equal(t, expectedData, data)
	})

	t.Run("error creating request", func(t *testing.T) {
		// Simulate an error while creating the request
		sync.url = "invalid-url"

		// Call the loadHTTPData method
		_, err := sync.loadHTTPData()

		// Assert that the expected error occurred
		assert.Error(t, err)
		assert.Equal(t, "error creating request", err.Error())
	})

	t.Run("error getting exchange rates", func(t *testing.T) {
		// Simulate an error while making the HTTP request
		server.Close()

		// Call the loadHTTPData method
		_, err := sync.loadHTTPData()

		// Assert that the expected error occurred
		assert.Error(t, err)
		assert.Equal(t, "error getting exchange rates", err.Error())
	})

	t.Run("error decoding exchange rates", func(t *testing.T) {
		// Create a test server that returns an invalid XML response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("invalid-xml"))
		}))
		defer server.Close()

		// Update the sync instance with the new test server URL
		sync.url = server.URL

		// Call the loadHTTPData method
		_, err := sync.loadHTTPData()

		// Assert that the expected error occurred
		assert.Error(t, err)
		assert.Equal(t, "error decoding exchange rates", err.Error())
	})
}
