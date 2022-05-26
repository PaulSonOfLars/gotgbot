package ext

import (
	"encoding/json"
	"net/url"
	"testing"
)

func TestValidateLoginQuery(t *testing.T) {
	// Make a sample query with known-good values.
	// id, first_name, last_name, username, photo_url, auth_date and hash fields
	query := url.Values{}
	query.Set("id", "12345")
	query.Set("first_name", "John")
	query.Set("last_name", "Smith")
	query.Set("username", "MrSmith")
	query.Set("photo_url", "example.com")
	query.Set("auth_date", "12345")
	query.Set("hash", "67fbf533a5e28a9bf93e12ca06b706e91383b96f51b9486af229ddcbc07ea801")

	t.Run("valid", func(t *testing.T) {
		ok, err := ValidateLoginQuery(query, "test_token")
		if err != nil {
			t.Errorf("failed to validate login query: %v", err)
			return
		}
		if !ok {
			t.Errorf("ValidateLoginQuery() with valid values should be true")
		}
	})
	t.Run("invalid", func(t *testing.T) {
		ok, err := ValidateLoginQuery(query, "invalid_token")
		if err != nil {
			t.Errorf("failed to validate login query: %v", err)
			return
		}
		if ok {
			t.Errorf("ValidateLoginQuery() with invalid values should be false")
		}
	})

	// If no hash is provided, fail
	query.Del("hash")
	t.Run("no hash", func(t *testing.T) {
		ok, err := ValidateLoginQuery(query, "invalid_token")
		if err != nil {
			t.Errorf("failed to validate login query: %v", err)
			return
		}
		if ok {
			t.Errorf("ValidateLoginQuery() with no hash values should be false")
		}
	})
}

func TestValidateWebApp(t *testing.T) {
	testUser, err := json.Marshal(struct {
		Id           int64  `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
		PhotoUrl     string `json:"photo_url"`
	}{
		Id:           12345,
		IsBot:        false,
		FirstName:    "John",
		LastName:     "Smith",
		Username:     "MrSmith",
		LanguageCode: "en-gb",
		PhotoUrl:     "example.com",
	})
	if err != nil {
		t.Error("Could not generate test user")
		return
	}
	// auth_date, query_id, user, and hash fields
	webAppQuery := url.Values{}
	webAppQuery.Set("auth_date", "12345")
	webAppQuery.Set("query_id", "12345")
	webAppQuery.Set("user", string(testUser))
	webAppQuery.Set("hash", "74e0df8d49d23ce7cb313f7baa2ee3fd96e51024aa398ad25ed6a90f91969746")

	t.Run("valid initData", func(t *testing.T) {
		ok, err := ValidateWebAppInitData(webAppQuery.Encode(), "test_token")
		if err != nil {
			t.Errorf("Failed to validate webapp query: %v", err)
			return
		}
		if !ok {
			t.Errorf("ValidateWebAppInitData() with valid values should be true")
		}
	})
	t.Run("invalid initData", func(t *testing.T) {
		ok, err := ValidateWebAppInitData(webAppQuery.Encode(), "invalid_token")
		if err != nil {
			t.Errorf("Failed to validate webapp query: %v", err)
			return
		}
		if ok {
			t.Errorf("ValidateWebAppInitData() with invalid values should be false")
		}
	})

	t.Run("valid query", func(t *testing.T) {
		ok, err := ValidateWebAppQuery(webAppQuery, "test_token")
		if err != nil {
			t.Errorf("Failed to validate webapp query: %v", err)
			return
		}
		if !ok {
			t.Errorf("ValidateWebAppQuery() with valid values should be true")
		}
	})
	t.Run("invalid query", func(t *testing.T) {
		ok, err := ValidateWebAppQuery(webAppQuery, "invalid_token")
		if err != nil {
			t.Errorf("Failed to validate webapp query: %v", err)
			return
		}
		if ok {
			t.Errorf("ValidateWebAppQuery() with invalid values should be false")
		}
	})
}
