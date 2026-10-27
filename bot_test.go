package gotgbot

import (
	"slices"
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
		bot      *Bot
		username string
		want     bool
	}{
		{
			name:     "nil bot",
			bot:      nil,
			username: "MyBot",
			want:     false,
		},
		{
			name:     "empty username",
			bot:      b,
			username: "",
			want:     false,
		},
		{
			name:     "at sign only",
			bot:      b,
			username: "@",
			want:     false,
		},
		{
			name:     "primary exact match",
			bot:      b,
			username: "MyBot",
			want:     true,
		},
		{
			name:     "primary case insensitive",
			bot:      b,
			username: "mybot",
			want:     true,
		},
		{
			name:     "primary with @ prefix",
			bot:      b,
			username: "@MyBot",
			want:     true,
		},
		{
			name:     "primary with @ prefix lowercase",
			bot:      b,
			username: "@mybot",
			want:     true,
		},
		{
			name:     "secondary exact match",
			bot:      b,
			username: "secondaryBot",
			want:     true,
		},
		{
			name:     "secondary case insensitive",
			bot:      b,
			username: "secondarybot",
			want:     true,
		},
		{
			name:     "secondary with @ prefix",
			bot:      b,
			username: "@secondaryBot",
			want:     true,
		},
		{
			name:     "third username defined with @ prefix matched without @",
			bot:      b,
			username: "thirdBot",
			want:     true,
		},
		{
			name:     "third username defined with @ prefix matched with @",
			bot:      b,
			username: "@thirdBot",
			want:     true,
		},
		{
			name:     "unknown bot",
			bot:      b,
			username: "otherBot",
			want:     false,
		},
		{
			name:     "unknown bot with @ prefix",
			bot:      b,
			username: "@otherBot",
			want:     false,
		},
		{
			name: "primary bot has @ in username field",
			bot: &Bot{
				User: User{Username: "@AtBot"},
			},
			username: "AtBot",
			want:     true,
		},
		{
			name: "missing username from disabled token check does not match missing string",
			bot: &Bot{
				User: User{Username: "<missing>"},
			},
			username: "<missing>",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bot.MatchesUsername(tt.username); got != tt.want {
				t.Errorf("MatchesUsername(%q) = %v, want %v", tt.username, got, tt.want)
			}
		})
	}
}

func TestBot_UpdateAdditionalUsernames(t *testing.T) {
	var nilBot *Bot
	// Should not panic on nil receiver
	nilBot.UpdateAdditionalUsernames([]string{"test"})
	if nilBot.AdditionalUsernames() != nil {
		t.Errorf("expected nil for nil bot AdditionalUsernames")
	}

	b := &Bot{}
	if got := b.AdditionalUsernames(); got != nil {
		t.Errorf("expected nil for uninitialized AdditionalUsernames, got %v", got)
	}

	b.UpdateAdditionalUsernames([]string{"@Alpha", "Beta", "", "@"})
	got := b.AdditionalUsernames()
	slices.Sort(got)
	expected := []string{"alpha", "beta"}
	if !slices.Equal(got, expected) {
		t.Errorf("AdditionalUsernames() = %v, want %v", got, expected)
	}

	// Update / replace
	b.UpdateAdditionalUsernames([]string{"gamma"})
	got = b.AdditionalUsernames()
	if len(got) != 1 || got[0] != "gamma" {
		t.Errorf("AdditionalUsernames() after replace = %v, want [gamma]", got)
	}
}

func TestBot_ConcurrentAccess(t *testing.T) {
	b := &Bot{
		User: User{Username: "PrimaryBot"},
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func(idx int) {
			defer wg.Done()
			b.UpdateAdditionalUsernames([]string{"alt1", "alt2", "alt3"})
		}(i)
		go func() {
			defer wg.Done()
			_ = b.MatchesUsername("alt1")
			_ = b.MatchesUsername("@PrimaryBot")
			_ = b.AdditionalUsernames()
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
	if !b.MatchesUsername("@ThirdBot") {
		t.Errorf("expected @ThirdBot to match")
	}
}
