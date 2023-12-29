package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	schemaURL          = "https://raw.githubusercontent.com/PaulSonOfLars/telegram-bot-api-spec/%s/api.json"
	latestCommitURL    = "https://api.github.com/repos/PaulSonOfLars/telegram-bot-api-spec/commits?page=1&per_page=1"
	specCommitFileName = "spec_commit"
)

func main() {
	// If GOTGBOT_UPGRADE set, get latest commit and generate from that.
	// If GOTGBOT_FROM_FILE is set, read from file instead of HTTP (allows offline dev)
	apiSpec, err := getAPISpec(os.Getenv("GOTGBOT_UPGRADE") != "", os.Getenv("GOTGBOT_FROM_FILE"))
	if err != nil {
		panic(fmt.Errorf("failed to get API spec: %w", err))
	}

	err = generate(apiSpec)
	if err != nil {
		panic(fmt.Errorf("failed to generate telegram bot api library from latest API spec: %w", err))
	}
}

func getCommit(upgrade bool) (string, error) {
	if upgrade {
		commit, err := updatePinnedCommit()
		if err != nil {
			return "", fmt.Errorf("failed to update pinned commit: %w", err)
		}
		fmt.Printf("Generating library from latest commit %s\n", commit) // nolint
		return commit, nil
	}

	// Else, use the pinned commit, for reproducible builds.
	contents, err := os.ReadFile(specCommitFileName)
	if err != nil {
		return "", fmt.Errorf("failed to read file spec: %w", err)
	}

	commit := strings.TrimSpace(string(contents))
	fmt.Printf("Generating library from pinned commit %s\n", commit) // nolint
	return commit, nil
}

func updatePinnedCommit() (string, error) {
	type APIResponse struct {
		Sha string `json:"sha"`
	}

	// Get the latest commit from github
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, latestCommitURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for latest commit: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get latest commit: %w", err)
	}
	defer resp.Body.Close()

	var res []APIResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("failed to decode GET request to update commit: %w", err)
	}

	commit := res[0].Sha

	err = os.WriteFile(specCommitFileName, []byte(commit), 0600)
	if err != nil {
		return "", fmt.Errorf("failed to update commit pin file: %w", err)
	}

	return commit, nil
}

func getAPISpec(shouldUpgrade bool, fromFile string) (APIDescription, error) {
	if fromFile != "" {
		if shouldUpgrade {
			return APIDescription{}, errors.New("upgrade and from_file are mutually exclusive options")
		}

		bs, err := os.ReadFile(fromFile)
		if err != nil {
			return APIDescription{}, fmt.Errorf("failed to read file: %w", err)
		}

		var d APIDescription
		if err := json.NewDecoder(bytes.NewReader(bs)).Decode(&d); err != nil {
			return APIDescription{}, fmt.Errorf("failed to decode API JSON: %w", err)
		}
		return d, nil
	}

	commit, err := getCommit(shouldUpgrade)
	if err != nil {
		panic(fmt.Errorf("failed to get commit to use for spec: %w", err))
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf(schemaURL, commit), nil)
	if err != nil {
		return APIDescription{}, fmt.Errorf("failed to create request for telegram bot api spec: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return APIDescription{}, fmt.Errorf("failed to download telegram bot api spec at %s: %w", commit, err)
	}
	defer resp.Body.Close()

	var d APIDescription
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return APIDescription{}, fmt.Errorf("failed to decode API JSON: %w", err)
	}
	return d, nil
}
