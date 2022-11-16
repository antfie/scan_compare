package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
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
		panic(err)
	}

	var credentialsFilePath = filepath.Join(homePath, ".veracode", "credentials")

	if _, err := os.Stat(credentialsFilePath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("Could not find a Veracode credentials file. See: https://docs.veracode.com/r/c_configure_api_cred_file")
	}

	file, err := os.Open(credentialsFilePath)

	if err != nil {
		panic("Could not open the Veracode credentials file. See: https://docs.veracode.com/r/c_configure_api_cred_file")
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

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	panic("Could not process the Veracode credentials file. See: https://docs.veracode.com/r/c_configure_api_cred_file")
}
