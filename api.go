package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/antfie/veracode-go-hmac-authentication/hmac"
	"github.com/fatih/color"
)

type API struct {
	id     string
	key    string
	region string
}

func (api API) makeApiRequest(apiUrl, httpMethod string) []byte {
	if api.region == "us" {
		apiUrl = strings.Replace(apiUrl, ".com", ".us", 1)
	} else if api.region == "eu" {
		apiUrl = strings.Replace(apiUrl, ".com", ".eu", 1)
	}

	parsedUrl, err := url.Parse(apiUrl)

	if err != nil {
		color.Red("Error: Invalid API URL")
		os.Exit(1)
	}

	client := &http.Client{}
	req, err := http.NewRequest(httpMethod, parsedUrl.String(), nil)

	if err != nil {
		color.Red("Error: Could not create API request")
		os.Exit(1)
	}

	authorizationHeader, err := hmac.CalculateAuthorizationHeader(parsedUrl, httpMethod, api.id, api.key)

	if err != nil {
		color.Red("Error: Could not calculate the authorization header")
		os.Exit(1)
	}

	req.Header.Add("Authorization", authorizationHeader)
	req.Header.Add("User-Agent", fmt.Sprintf("ScanCompare/%s", AppVersion))

	resp, err := client.Do(req)

	if err != nil {
		color.Red("Error: There was a problem communicating with the API. Please check your connectivity and the service status page at https://status.veracode.com")
		os.Exit(1)
	}

	if resp.StatusCode == 401 {
		color.Red("Error: You are not authorized to perform this action. Please check your credentials are valid for this Veracode region and that you have the correct permissions. For help contact your Veracode administrator.")
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		color.Red(fmt.Sprintf("Error: API request returned status of %s", resp.Status))
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		color.Red("Error: There was a problem processing the API response. Please check your connectivity and the service status page at https://status.veracode.com")
		os.Exit(1)
	}

	return body
}
