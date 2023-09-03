package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/tlkamp/litter-api/v2/internal/auth"

	"github.com/pkg/errors"
)

const (
	apiURL = "https://v2.api.whisker.iothings.site"
	apiKey = "p7ndMoj61npRZP5CVz9v4Uj0bG769xy6758QRBPb"
)

// Client is the LitterRobot API client.
type Client struct {
	mutex  *sync.RWMutex
	auth   *auth.Client
	api    *http.Client
	robots map[string]Robot
}

// New returns an initialized *Client.
func New(email, password string) *Client {
	auth := auth.New(email, password)
	return &Client{
		mutex:  &sync.RWMutex{},
		auth:   auth,
		api:    http.DefaultClient,
		robots: make(map[string]Robot),
	}
}

// Login authenticates the client.
func (c *Client) Login(ctx context.Context) error {
	return c.auth.Login(ctx)
}

// RefreshToken can be used to periodically refresh the access token created during the login process.
// Use this function in a goroutine for long-running programs.
func (c *Client) RefreshToken(ctx context.Context) error {
	return c.auth.DoRefreshToken(ctx)
}

// FetchRobots fetches the robots from the LitterRobot API.
// The robots are cached on the client and can be fetched without additional network calls using Robots() or Robot(id)
func (c *Client) FetchRobots(ctx context.Context) error {
	robotEndpoint := fmt.Sprintf("/users/%s/robots", c.auth.UserID())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL+robotEndpoint, nil)
	if err != nil {
		return errors.Wrap(err, "error creating robot request")
	}

	req.Header = map[string][]string{
		"Authorization": []string{fmt.Sprintf("Bearer %s", c.auth.IDToken())},
		"x-api-key":     []string{apiKey},
	}

	resp, err := c.api.Do(req)
	if err != nil {
		return errors.Wrap(err, "error fetching robots")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bd, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading robots body")
	}

	var robots []robotResponse

	if err := json.Unmarshal(bd, &robots); err != nil {
		return errors.Wrap(err, "error unmarshaling robots response")
	}

	r := make(map[string]Robot, len(robots))
	for _, rb := range robots {
		rbt := newRobot(rb)
		r[rbt.LitterRobotID] = rbt
	}

	c.mutex.Lock()
	c.robots = r
	c.mutex.Unlock()

	return nil
}

// Robot returns a fetched robot with the corresponding ID.
func (c *Client) Robot(id string) Robot {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.robots[id]
}

// Robots returns all fetched robots.
func (c *Client) Robots() []Robot {
	robots := make([]Robot, 0, len(c.robots))

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, r := range c.robots {
		robots = append(robots, r)
	}
	return robots
}

// FetchInsights returns the Litter Robot insights for the specified period. This function always makes a network call.
func (c *Client) FetchInsights(ctx context.Context, id string, days, tzOffset int) (*Insight, error) {
	insightsURL := fmt.Sprintf("/users/%s/robots/%s/insights", c.auth.UserID(), id)

	if days < 1 {
		return nil, errors.New("days must be greather than 0")
	}

	params := url.Values{}
	params.Set("days", fmt.Sprintf("%d", days))
	params.Set("timezoneOffset", fmt.Sprintf("%d", tzOffset))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL+insightsURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating insights request")
	}

	req.Header = map[string][]string{
		"Authorization": []string{fmt.Sprintf("Bearer %s", c.auth.IDToken())},
		"x-api-key":     []string{apiKey},
	}

	resp, err := c.api.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error sending insights request")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading insights body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var insight *Insight
	if err := json.Unmarshal(body, &insight); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling insights response")
	}

	return insight, nil
}

func (c *Client) sendCommand(ctx context.Context, robotId, command string) error {
	cmdPath := fmt.Sprintf("/users/%s/robots/%s/dispatch-commands", c.auth.UserID(), robotId)

	cmd := &commandBody{
		Command:       command,
		LitterRobotId: robotId,
	}

	body, err := json.Marshal(cmd)
	if err != nil {
		return errors.Wrap(err, "error marshaling command body")
	}

	payload := bytes.NewBuffer(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL+cmdPath, payload)
	if err != nil {
		return errors.Wrap(err, "error building command request")
	}

	req.Header = map[string][]string{
		"Authorization": []string{fmt.Sprintf("Bearer %s", c.auth.IDToken())},
		"x-api-key":     []string{apiKey},
	}

	resp, err := c.api.Do(req)
	if err != nil {
		return errors.Wrap(err, "error sending command")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// PowerOn - Turn unit power on.
func (c *Client) PowerOn(ctx context.Context, robotId string) error {
	return c.sendCommand(ctx, robotId, powerCmd+on)
}

// PowerOff - Turn unit power off.
func (c *Client) PowerOff(ctx context.Context, robotId string) error {
	return c.sendCommand(ctx, robotId, powerCmd+off)
}

// NightLightOn - Turn nightlight on.
func (c *Client) NightLightOn(ctx context.Context, robotId string) error {
	return c.sendCommand(ctx, robotId, nightLightCmd+on)
}

// NightLightOff - Turn nightlight off.
func (c *Client) NightLightOff(ctx context.Context, robotId string) error {
	return c.sendCommand(ctx, robotId, nightLightCmd+off)
}

// PanelLockOn - Enable the panel lock.
func (c *Client) PanelLockOn(ctx context.Context, robotId string) error {
	return c.sendCommand(ctx, robotId, panelLockCmd+on)
}

// PanelLockOff - Disable the panel lock.
func (c *Client) PanelLockOff(ctx context.Context, robotId string) error {
	return c.sendCommand(ctx, robotId, panelLockCmd+off)
}

// Cycle - Start a clean cycle.
func (c *Client) Cycle(ctx context.Context, robotId string) error {
	return c.sendCommand(ctx, robotId, cycleCmd)
}

// Wait - Set clean cycle wait time.
func (c *Client) Wait(ctx context.Context, robotId string, val string) error {
	return c.sendCommand(ctx, robotId, waitCmd+val)
}
