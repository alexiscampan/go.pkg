// Package http to define utils function for http
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

const (
	// ConnectMaxWaitTime to have a connection limit time
	ConnectMaxWaitTime = 20 * time.Second
)

// GetQuery to send Get http request
func GetQuery(ctx context.Context, url string, params, headers map[string]string) (resp []byte, err error) {
	client := initClient()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.URL.RawQuery = q.Encode()
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request %q: %w", url, err)
	}

	defer func() {
		if err2 := response.Body.Close(); err2 != nil {
			if err == nil {
				err = fmt.Errorf("failed to close response from %q: %w", url, err2)
			}
		}
	}()
	startRead := time.Now()
	resp, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from %q: %w", url, err)
	}
	endRead := time.Now()
	fmt.Printf("Read response took: %s", endRead.Sub(startRead).String())
	return resp, err
}

// PostQuery to send Post http request
func PostQuery(ctx context.Context, url string, body []byte, params, headers map[string]string) (respCode int, err error) {
	client := initClient()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return http.StatusBadRequest, err
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.URL.RawQuery = q.Encode()
	response, err := client.Do(req)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to post %q: %w", url, err)
	}

	defer func() {
		if err2 := response.Body.Close(); err2 != nil {
			if err == nil {
				err = fmt.Errorf("failed to close response from %q: %w", url, err2)
			}
		}
	}()

	var result map[string]interface{}
	derr := json.NewDecoder(response.Body).Decode(&result)
	if derr != nil {
		return http.StatusForbidden, fmt.Errorf("failed to decode result %s: %w", response.Body, err)
	}
	if response.StatusCode != http.StatusCreated {
		return response.StatusCode, fmt.Errorf("failed to decode result %s: %w", response.Body, err)
	}
	return response.StatusCode, nil
}

func initClient() http.Client {
	return http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: ConnectMaxWaitTime,
			}).DialContext,
		},
	}
}
