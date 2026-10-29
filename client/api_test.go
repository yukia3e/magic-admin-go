package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yukia3e/magic-admin-go"
)

const testSecret = "sk_test_E123E4567E8901D2"

var testClientInfo = magic.ClientInfo{
	ClientId: "did:magic:731848cc-084e-41ff-bbdf-7f103817ea6b",
	AppScope: "global",
}

func TestNew(t *testing.T) {
	srv := createServer(t)
	defer srv.Close()

	// Replace host url to test one.
	client := magic.NewDefaultClient()
	client.SetBaseURL(srv.URL)

	api, err := New(testSecret, client)
	require.NoError(t, err, "can't create new client")

	assert.Equal(t, testClientInfo, api.ClientInfo)
	assert.NotNil(t, api.User)
}

func TestNewWrongSecret(t *testing.T) {
	srv := createServer(t)
	defer srv.Close()

	// Replace host url to test one.
	client := magic.NewDefaultClient()
	client.SetBaseURL(srv.URL)

	api, err := New("wrong_secret", client)
	require.Error(t, err, "server must return error")
	assert.Nil(t, api)
	_, ok := err.(*magic.AuthenticationError)
	require.True(t, ok, "Error type must be AuthenticationError")
}

func TestNewConnectionFailure(t *testing.T) {
	srv := createServer(t)
	srv.Close()

	// Use client without retries to fail fast.
	client := magic.NewClient()
	client.SetBaseURL(srv.URL)

	api, err := New(testSecret, client)
	require.Error(t, err, "request must fail")
	assert.Nil(t, api)
	_, ok := err.(*magic.APIConnectionError)
	require.True(t, ok, "Error type must be APIConnectionError")
}

// createServer creates internal server which simulates backend api for client info requests.
func createServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Method: %v", r.Method)
		t.Logf("Path: %v", r.URL.Path)

		w.Header().Add("Content-Type", "application/json")

		if r.Header.Get(magic.APISecretHeader) != testSecret {
			w.WriteHeader(http.StatusUnauthorized)
			resp := magic.Response{
				ErrorCode: "err_code_unauthorized",
				Message:   "unauthorized",
				Status:    "fail",
			}
			data, err := json.Marshal(resp)
			assert.NoError(t, err, "can't marshal test data")
			_, _ = w.Write(data)
			return
		}

		if r.Method == http.MethodGet && r.URL.Path == clientInfoV1 {
			resp := magic.Response{
				Data:   testClientInfo,
				Status: "ok",
			}
			data, err := json.Marshal(resp)
			assert.NoError(t, err, "can't marshal test data")
			_, _ = w.Write(data)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
}
