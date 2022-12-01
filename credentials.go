package main

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func formatCredential(credential string) string {
	var parts = strings.Split(credential, "-")
	if len(parts) == 2 {
		return parts[1]
	}

	return credential
}

func getCredentials(id, key string) (string, string) {
	id = formatCredential(id)
	key = formatCredential(key)

	// First try CLI flags
	if len(id) == 32 && len(key) == 128 {
		return id, key
	}

	if len(id) > 0 && len(id) != 32 {
		color.HiRed("Error: Invalid value for -vid")
		os.Exit(1)
	}

	if len(key) > 0 && len(key) != 128 {
		color.HiRed("Error: Invalid value for -vkey")
		os.Exit(1)
	}

	if len(id) > 0 && len(key) == 0 || len(key) > 0 && len(id) == 0 {
		color.HiRed("Error: If passing Veracode API key via command line both -vid and -vkey are required")
		os.Exit(1)
	}

	id = ""
	key = ""

	// Then try environment variables
	id = os.Getenv("VERACODE_API_KEY_ID")
	key = os.Getenv("VERACODE_API_KEY_SECRET")

	id = formatCredential(id)
	key = formatCredential(key)

	if len(id) == 32 && len(key) == 128 {
		return id, key
	}

	if len(id) > 0 && len(id) != 32 {
		color.HiRed("Error: Invalid value for VERACODE_API_KEY_ID")
		os.Exit(1)
	}

	if len(key) > 0 && len(key) != 128 {
		color.HiRed("Error: Invalid value for VERACODE_API_KEY_SECRET")
		os.Exit(1)
	}

	if len(id) > 0 && len(key) == 0 || len(key) > 0 && len(id) == 0 {
		color.HiRed("Error: If passing Veracode API key via environment variables both VERACODE_API_KEY_ID and VERACODE_API_KEY_SECRET are required")
		os.Exit(1)
	}

	id = ""
	key = ""

	// Finally look for a Veracode credentials file
	homePath, err := os.UserHomeDir()

	if err != nil {
		color.HiRed("Error: Could not locate your home directory")
		os.Exit(1)
	}

	var credentialsFilePath = filepath.Join(homePath, ".veracode", "credentials")

	if _, err := os.Stat(credentialsFilePath); errors.Is(err, os.ErrNotExist) {
		color.HiRed("Error: Could not resolve any API credentials. Use either -vid and -vkey command line arguments, set VERACODE_API_KEY_ID and VERACODE_API_KEY_SECRET environment variables or create a Veracode credentials file - see: https://docs.veracode.com/r/c_configure_api_cred_file")
		os.Exit(1)
	}

	file, err := os.Open(credentialsFilePath)

	if err != nil {
		color.HiRed("Error: Could not open the Veracode credentials file. See: https://docs.veracode.com/r/c_configure_api_cred_file")
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
			if strings.Contains(line, "[") {
				found = false
				id = ""
				key = ""
			} else {
				if strings.Contains(line, "veracode_api_key_id = ") {
					id = strings.TrimSpace(strings.Split(line, "veracode_api_key_id = ")[1])
				}

				if strings.Contains(line, "veracode_api_key_secret = ") {
					key = strings.TrimSpace(strings.Split(line, "veracode_api_key_secret = ")[1])
				}

				if len(id) > 0 && len(key) > 0 {
					id = formatCredential(id)
					key = formatCredential(key)

					if len(id) != 32 {
						color.HiRed("Error: Invalid value for veracode_api_key_id in file \"%s\"", credentialsFilePath)
						os.Exit(1)
					}

					if len(key) != 128 {
						color.HiRed("Error: Invalid value for veracode_api_key_secret in file \"%s\"", credentialsFilePath)
						os.Exit(1)
					}

					return id, key

				}
			}
		}
	}

	color.HiRed("Error: Could not parse credentials from the Veracode credentials file. See: https://docs.veracode.com/r/c_configure_api_cred_file")
	os.Exit(1)
	return "", ""
}
