package http_client

import (
	"ModeFlex/data"
	"io"
	"log"
	"net/http"
	"time"
)

func clientTokenCheck() {
	expired := data.ClientTokenExpires < time.Now().UTC().Unix()
	if expired {
		e := ClientCredentialsGrant()
		if e != nil {
			log.Fatalf("Failed to refresh token! %v", e)
		}
	}
}

func APIRequestWithHeaders(method string, url string, body io.Reader) (*http.Request, error) {
	clientTokenCheck()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+data.ClientToken)
	return req, nil
}
