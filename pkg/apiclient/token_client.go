package apiclient

import (
	"encoding/json"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Sub         string `json:"sub"`
}

type ITokenApiClient interface {
	GetClientAccessToken() (string, error)
}

type TokenApiClient struct {
	TokenEndpoint string
	ClientId      string
	ClientSecret  string
}

func NewTokenApiClient() ITokenApiClient {
	return &TokenApiClient{
		TokenEndpoint: os.Getenv("TOKEN_ENDPOINT"),
		ClientId:      os.Getenv("CLIENT_ID"),
		ClientSecret:  os.Getenv("CLIENT_SECRET"),
	}
}

func (t *TokenApiClient) GetClientAccessToken() (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type": "client_credentials",
		}).
		SetBasicAuth(t.ClientId, t.ClientSecret).
		Post(t.TokenEndpoint)

	if err != nil {
		return "", errors.Errorf("Failed to call token endpoint: %v", err)
	}

	if resp.StatusCode() != 200 {
		return "", errors.Errorf("Token endpoint request failed with status code %d", resp.StatusCode())
	}

	tokenResponse := TokenResponse{}
	err = json.Unmarshal(resp.Body(), &tokenResponse)
	if err != nil {
		return "", errors.Wrap(err, "Failed to unmarshal token response")
	}

	return tokenResponse.AccessToken, nil
}
