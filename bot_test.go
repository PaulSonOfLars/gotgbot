package gotgbot

import (
	"sync"
	"testing"
)

func TestBot_MatchesUsername(t *testing.T) {
	b := &Bot{
		User: User{Username: "MyBot"},
	}
	b.UpdateAdditionalUsernames([]string{"secondaryBot", "@thirdBot"})

	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{
			name:     "empty username",
			username: "",
			want:     false,
		},
		{
			name:     "primary exact match",
			username: "MyBot",
			want:     true,
		},
		{
			name:     "primary case insensitive",
			username: "mybot",
			want:     true,
		},
		{
			name:     "secondary exact match",
			username: "secondaryBot",
			want:     true,
		},
		{
			name:     "secondary case insensitive",
			username: "secondarybot",
			want:     true,
		},
		{
			name:     "third username defined with @ prefix matched without @",
			username: "thirdBot",
			want:     true,
		},
		{
			name:     "unknown bot",
			username: "otherBot",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := b.MatchesUsername(tt.username); got != tt.want {
				t.Errorf("MatchesUsername(%q) = %v, want %v", tt.username, got, tt.want)
			}
		})
	}
}

func TestBot_ConcurrentAccess(t *testing.T) {
	b := &Bot{
		User: User{Username: "PrimaryBot"},
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			b.UpdateAdditionalUsernames([]string{"alt1", "alt2", "alt3"})
		}()
		go func() {
			defer wg.Done()
			_ = b.MatchesUsername("alt1")
			_ = b.MatchesUsername("PrimaryBot")
		}()
	}
	wg.Wait()
}

func TestNewBot_WithAdditionalUsernames(t *testing.T) {
	b, err := NewBot("123:abc", &BotOpts{
		DisableTokenCheck:   true,
		AdditionalUsernames: []string{"@SecondBot", "ThirdBot"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !b.MatchesUsername("secondbot") {
		t.Errorf("expected secondbot to match")
	}
	if !b.MatchesUsername("thirdbot") {
		t.Errorf("expected thirdbot to match")
	}
}
