package main

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func getCredentials(id, key string) (string, string) {
	// First try CLI flags
	if len(id) == 32 && len(key) == 128 {
		return id, key
	}

	// Then try environment variables
	id = os.Getenv("VERACODE_API_KEY_ID")
	key = os.Getenv("VERACODE_API_KEY_SECRET")

	if len(id) == 32 && len(key) == 128 {
		return id, key
	}

	// Finally look for a Veracode credentials file
	homePath, err := os.UserHomeDir()

	if err != nil {
		color.Red("Error: Could not locate your home directory")
		os.Exit(1)
	}

	var credentialsFilePath = filepath.Join(homePath, ".veracode", "credentials")

	if _, err := os.Stat(credentialsFilePath); errors.Is(err, os.ErrNotExist) {
		color.Red("Error: Could not find a Veracode credentials file. See: https://docs.veracode.com/r/c_configure_api_cred_file")
		os.Exit(1)
	}

	file, err := os.Open(credentialsFilePath)

	if err != nil {
		color.Red("Error: Could not open the Veracode credentials file. See: https://docs.veracode.com/r/c_configure_api_cred_file")
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	found := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "[default]" {
			found = true
		} else if found {
			if strings.Contains(line, "veracode_api_key_id = ") {
				id = strings.Split(line, "veracode_api_key_id = ")[1]
			}

			if strings.Contains(line, "veracode_api_key_secret = ") {
				key = strings.Split(line, "veracode_api_key_secret = ")[1]
			}
		}

		if len(id) == 32 && len(key) == 128 {
			return id, key
		}
	}

	color.Red("Error: Could not process the Veracode credentials file. See: https://docs.veracode.com/r/c_configure_api_cred_file")
	os.Exit(1)
	return "", ""
}
