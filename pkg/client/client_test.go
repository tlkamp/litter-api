package client

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	a := assert.New(t)
	email := os.Getenv("LR_EMAIL")
	password := os.Getenv("LR_PASSWORD")

	c := New(email, password)

	ctx := context.Background()

	a.NoError(c.Login(ctx))

	a.NoError(c.FetchRobots(ctx))
	a.NotNil(c.robots)

	id := c.Robots()[0].LitterRobotID

	in, err := c.FetchInsights(ctx, id, 4, -7)
	a.NoError(err)
	a.NotZero(in)
	a.NotZero(in.CycleHistory)
	a.Greater(len(in.CycleHistory), 0)

	a.NoError(c.Cycle(ctx, id))
}

func TestSetToken(t *testing.T) {
	c := New("", "")
	c.SetToken("testing")
	assert.Equal(t, "testing", c.Token())
}
