package http_client

import (
	"ModeFlex/api"
	"ModeFlex/data"
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func ClientCredentialsGrant() error {
	v, err := json.Marshal(&api.ClientCredentialsGrantPOSTBody{
		ClientID:     data.Configuration.ClientID,
		ClientSecret: data.Configuration.ClientSecret,
		GrantType:    "client_credentials",
		Scope:        "public",
	})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", api.OsuOAuthObtainToken, bytes.NewBuffer(v))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := data.OsuHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var clientCredentials *api.ClientCredentialsTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&clientCredentials)
	if err != nil {
		return err
	}

	data.ClientToken = clientCredentials.AccessToken
	data.ClientTokenExpires = time.Now().UTC().Unix() + clientCredentials.ExpiresIn
	return nil
}
