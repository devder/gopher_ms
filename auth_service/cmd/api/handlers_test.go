package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RoundTripFunc func(*http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestAuthenticate(t *testing.T) {
	jsonToReturn := `
		{
			"error": false,
			"message": "some message"
		}
	`

	// mock calling the broker
	client := NewTestClient(func(r *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(jsonToReturn)),
			Header:     make(http.Header),
		}
	})

	testApp.Client = client

	postBody := map[string]string{
		"email":    "me@ex.co",
		"password": "password",
	}

	body, err := json.Marshal(postBody)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/authenticate", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.Authenticate)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("Expected %d, but got %d", http.StatusAccepted, rec.Code)
	}
}
