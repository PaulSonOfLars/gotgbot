package ext

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

// ValidateLoginQuery validates a login widget query.
// See https://core.telegram.org/widgets/login#checking-authorization for more details.
func ValidateLoginQuery(query url.Values, token string) bool {
	return validateQuery(query, getSHA256(token))
}

// ValidateWebAppInitData validates a webapp's initData field for safe use on the server-side.
// The initData fiels is stored as a query string, so this is converted and then validated.
// See https://core.telegram.org/bots/webapps#validating-data-received-via-the-web-app for more details.
func ValidateWebAppInitData(initData string, token string) (bool, error) {
	query, err := url.ParseQuery(initData)
	if err != nil {
		return false, err
	}

	return ValidateWebAppQuery(query, token), nil
}

// ValidateWebAppQuery validates a webapp's initData query for safe use on the server side.
// The input is expected to be the parsed initData query string.
// See https://core.telegram.org/bots/webapps#validating-data-received-via-the-web-app for more details.
func ValidateWebAppQuery(query url.Values, token string) bool {
	return validateQuery(query, generateHMAC256(token, []byte("WebAppData")))
}

func validateQuery(query url.Values, secretKey []byte) bool {
	// If no hash, we can't check; fail-fast.
	hash := query.Get("hash")
	if hash == "" {
		return false
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
	expectedHMAC := generateHMAC256(dataCheck, secretKey)

	// Hex encode expected hmac_256 value.
	expectedHex := getHex(expectedHMAC)

	// Check hash matches as expected.
	return hmac.Equal(expectedHex, []byte(hash))
}

func generateHMAC256(data string, secretKey []byte) []byte {
	hmac256Writer := hmac.New(sha256.New, secretKey)
	hmac256Writer.Write([]byte(data))
	return hmac256Writer.Sum(nil)
}

func getSHA256(data string) []byte {
	sha := sha256.New()
	sha.Write([]byte(data))
	return sha.Sum(nil)
}

func getHex(expectedHMAC []byte) []byte {
	expectedHex := make([]byte, hex.EncodedLen(len(expectedHMAC)))
	hex.Encode(expectedHex, expectedHMAC)
	return expectedHex
}
