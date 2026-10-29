package magic

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorString(t *testing.T) {
	err := &Error{Response: Response{
		ErrorCode: "err_code_unauthorized",
		Message:   "unauthorized",
		Status:    "fail",
	}}

	assert.JSONEq(t,
		`{"data":null,"error_code":"err_code_unauthorized","message":"unauthorized","status":"fail"}`,
		err.Error())
}

func TestWrapError(t *testing.T) {
	apiErr := &Error{Response: Response{Message: "failure", Status: "fail"}}

	tests := []struct {
		name       string
		statusCode int
		want       error
	}{
		{name: "forbidden", statusCode: http.StatusForbidden, want: &ForbiddenError{Err: apiErr}},
		{name: "bad request", statusCode: http.StatusBadRequest, want: &BadRequestError{Err: apiErr}},
		{name: "too many requests", statusCode: http.StatusTooManyRequests, want: &RateLimitingError{Err: apiErr}},
		{name: "unauthorized", statusCode: http.StatusUnauthorized, want: &AuthenticationError{Err: apiErr}},
		{name: "internal server error", statusCode: http.StatusInternalServerError, want: &APIError{Err: apiErr}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &resty.Response{RawResponse: &http.Response{StatusCode: tt.statusCode}}

			err := WrapError(r, apiErr)
			assert.Equal(t, tt.want, err)
			assert.Equal(t, apiErr.Error(), err.Error())
		})
	}
}

func TestAPIConnectionError(t *testing.T) {
	err := &APIConnectionError{Err: errors.New("connection refused")}

	assert.Equal(t, "connection refused", err.Error())
}

func TestNewClient(t *testing.T) {
	client := NewClient()

	assert.Equal(t, APIURL, client.BaseURL)
}

func TestNewDefaultClient(t *testing.T) {
	client := NewDefaultClient()

	assert.Equal(t, APIURL, client.BaseURL)
	assert.Equal(t, 3, client.RetryCount)
	assert.Equal(t, time.Second, client.RetryWaitTime)
	assert.Equal(t, 10*time.Second, client.RetryMaxWaitTime)
}

func TestNewClientWithRetry(t *testing.T) {
	client := NewClientWithRetry(5, 2*time.Second, 30*time.Second)

	assert.Equal(t, APIURL, client.BaseURL)
	assert.Equal(t, 5, client.RetryCount)
	assert.Equal(t, 2*time.Second, client.RetryWaitTime)
	assert.Equal(t, 30*time.Second, client.RetryMaxWaitTime)
	require.NotNil(t, client.RetryAfter)

	wait, err := client.RetryAfter(client, nil)
	assert.Zero(t, wait)
	assert.ErrorIs(t, err, ErrRespQuotaExceeded)
}
