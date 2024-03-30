package sync

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExchangeRateSync_loadHTTPDataWithHTTPTest(t *testing.T) {
	// Example XML response
	responseXML := `<Envelope><Subject>Subject</Subject><Sender><name>Test Sender</name></Sender><Cube><Cube time="2023-01-01"><Cube currency="USD" rate="1.1"/><Cube currency="EUR" rate="1.2"/></Cube></Cube></Envelope>`

	tests := []struct {
		name             string
		setupHandler     func(t *testing.T) http.Handler
		expectErr        bool
		expectedRateSize int
	}{
		{
			name: "Successful fetch",
			setupHandler: func(t *testing.T) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(responseXML))
				})
			},
			expectErr:        false,
			expectedRateSize: 2,
		},
		{
			name: "Network error",
			setupHandler: func(t *testing.T) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup httptest server
			server := httptest.NewServer(tc.setupHandler(t))
			defer server.Close()

			ers := NewExchangeRateSync("public", server.URL, nil) // Use the test server URL

			result, err := ers.loadHTTPData()
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tc.expectedRateSize)
				if tc.expectedRateSize > 0 {
					assert.Equal(t, "USD", result[0].Currency)
					assert.Equal(t, 1.1, result[0].Rate)
				}
			}
		})
	}
}
