package main

import (
	"fmt"
	"net/http"
)

const schemaURL = "https://raw.githubusercontent.com/PaulSonOfLars/telegram-bot-api-spec/main/api.json"

func main() {
	// Get the latest bot API spec from github
	resp, err := http.Get(schemaURL)
	if err != nil {
		panic(fmt.Errorf("failed to download telegram bot api spec at '%s': %w", schemaURL, err))
	}
	defer resp.Body.Close()

	err = generate(resp.Body)
	if err != nil {
		panic(fmt.Errorf("failed to generate telegram bot api library from latest API spec: %w", err))
	}
}
