package auth

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	a := assert.New(t)

	email := os.Getenv("LR_EMAIL")
	password := os.Getenv("LR_PASSWORD")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("Login", func(t *testing.T) {
		c := New(email, password)

		a.NoError(c.Login(ctx))
		a.NotZero(t, c.userId)
	})

	t.Run("RefreshToken", func(t *testing.T) {
		c := New(email, password)
		a.NoError(c.Login(ctx))
		a.NoError(c.DoRefreshToken(ctx))
	})
}
