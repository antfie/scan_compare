package main

import (
	"io"
	"net/http"
	"net/url"

	"github.com/antfie/veracode-go-hmac-authentication/hmac"
)

func (api API) makeApiRequest(apiUrl, httpMethod string) []byte {
	parsedUrl, err := url.Parse(apiUrl)

	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(httpMethod, parsedUrl.String(), nil)

	if err != nil {
		panic(err)
	}

	authorizationHeader, err := hmac.CalculateAuthorizationHeader(parsedUrl, httpMethod, api.id, api.key)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", authorizationHeader)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic("Expected status 200. Status was: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return body
}
