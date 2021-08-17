package main

import (
	"context"
	"fmt"
	"net/http"
)

const schemaURL = "https://raw.githubusercontent.com/PaulSonOfLars/telegram-bot-api-spec/main/api.json"

func main() {
	// Get the latest bot API spec from github
	req, err := http.NewRequestWithContext(context.Background(), "GET", schemaURL, nil)
	if err != nil {
		panic(fmt.Errorf("failed to create request for telegram bot api spec: %w", err))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("failed to download telegram bot api spec: %w", err))
	}
	defer resp.Body.Close()

	err = generate(resp.Body)
	if err != nil {
		panic(fmt.Errorf("failed to generate telegram bot api library from latest API spec: %w", err))
	}
}
