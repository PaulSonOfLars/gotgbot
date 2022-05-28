package ext

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// ValidateLoginQuery validates a login widget query.
// See https://core.telegram.org/widgets/login#checking-authorization for more details.
func ValidateLoginQuery(query url.Values, token string) (bool, error) {
	tokenHash, err := getSHA256(token)
	if err != nil {
		return false, fmt.Errorf("failed to hash token: %w", err)
	}

	return validateQuery(query, tokenHash)
}

// ValidateWebAppInitData validates a webapp's initData field for safe use on the server-side.
// The initData field is stored as a query string, so this is converted and then validated.
// See https://core.telegram.org/bots/webapps#validating-data-received-via-the-web-app for more details.
func ValidateWebAppInitData(initData string, token string) (bool, error) {
	query, err := url.ParseQuery(initData)
	if err != nil {
		return false, fmt.Errorf("failed to parse URL query: %w", err)
	}

	return ValidateWebAppQuery(query, token)
}

// ValidateWebAppQuery validates a webapp's initData query for safe use on the server side.
// The input is expected to be the parsed initData query string.
// See https://core.telegram.org/bots/webapps#validating-data-received-via-the-web-app for more details.
func ValidateWebAppQuery(query url.Values, token string) (bool, error) {
	tokenHMAC, err := generateHMAC256(token, []byte("WebAppData"))
	if err != nil {
		return false, fmt.Errorf("failed to generate token HMAC: %w", err)
	}

	return validateQuery(query, tokenHMAC)
}

func validateQuery(query url.Values, secretKey []byte) (bool, error) {
	// If no hash, we can't check; fail-fast.
	hash := query.Get("hash")
	if hash == "" {
		// Should this be an error?
		return false, nil
	}

	// Make list of args for ordered sorting.
	// len()-1, because we ignore the hash key.
	args := make([]string, 0, len(query)-1)
	for x, y := range query {
		if x == "hash" {
			// ignore the hash
			continue
		}
		args = append(args, x+"="+y[0])
	}

	// Sort args to ensure consistency.
	sort.Strings(args)

	// Join data with newline, as defined by telegram.
	dataCheck := strings.Join(args, "\n")

	// Generate HMAC of expected data.
	expectedHMAC, err := generateHMAC256(dataCheck, secretKey)
	if err != nil {
		return false, fmt.Errorf("failed to generate data HMAC: %w", err)
	}

	// Hex encode expected hmac_256 value.
	expectedHex := getHex(expectedHMAC)

	// Check hash matches as expected.
	return hmac.Equal(expectedHex, []byte(hash)), nil
}

func generateHMAC256(data string, secretKey []byte) ([]byte, error) {
	hmac256Writer := hmac.New(sha256.New, secretKey)
	_, err := hmac256Writer.Write([]byte(data))
	if err != nil {
		return nil, err
	}

	return hmac256Writer.Sum(nil), nil
}

func getSHA256(data string) ([]byte, error) {
	sha256Writer := sha256.New()
	_, err := sha256Writer.Write([]byte(data))
	if err != nil {
		return nil, err
	}

	return sha256Writer.Sum(nil), nil
}

func getHex(expectedHMAC []byte) []byte {
	expectedHex := make([]byte, hex.EncodedLen(len(expectedHMAC)))
	hex.Encode(expectedHex, expectedHMAC)
	return expectedHex
}
