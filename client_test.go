package api

import (
	"os"
	"strings"
	"testing"
)

func TestClient_Integration(t *testing.T) {
	email, ok := os.LookupEnv("LR_EMAIL")
	if !ok {
		t.Fatal("No value for LR_EMAIL")
	}

	password, ok := os.LookupEnv("LR_PASSWORD")
	if !ok {
		t.Fatal("No value for LR_PASSWORD")
	}

	lc, err := NewClient(&Config{
		Email:        email,
		Password:     password,
		ApiKey:       "p7ndMoj61npRZP5CVz9v4Uj0bG769xy6758QRBPb",
		ClientSecret: "C63CLXOmwNaqLTB2xXo6QIWGwwBamcPuaul",
		ClientId:     "IYXzWN908psOm7sNpe4G.ios.whisker.robots",
	})
	if err != nil {
		t.Fatalf("Received error: %s", err)
	}

	t.Run("Test States", func(t *testing.T) {
		states, err := lc.States()
		if err != nil {
			t.Fatalf("Received error: %s", err)
		}

		if len(states) != 1 {
			t.Errorf("expected %d state(s), got: %d", 1, len(states))
		}

		for _, state := range states {
			if !strings.Contains(state.Name, "magic") {
				t.Fatalf("Robot state has invalid name: %s", state.Name)
			}
		}
	})
}
