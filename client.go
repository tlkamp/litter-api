package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

const (
	on            = "1"
	off           = "0"
	powerCmd      = "<P"
	cycleCmd      = "<C"
	nightLightCmd = "<N"
	panelLockCmd  = "<L"
	waitCmd       = "<W"
)

// NewClient - Initialize a new Client with provided configuration Config
func NewClient(cfg *Config) (*Client, error) {
	client := &Client{
		apiClient:  resty.New(),
		authClient: resty.New(),
		Config:     cfg,
		robots:     make(map[string]State),
	}

	if cfg.AuthUrl == "" {
		cfg.AuthUrl = "https://autopets.sso.iothings.site"
	}
	client.authClient.SetHeader("x-api-key", cfg.ApiKey)
	client.authClient.SetHostURL(cfg.AuthUrl)

	if cfg.ApiUrl == "" {
		cfg.ApiUrl = "https://v2.api.whisker.iothings.site"
	}

	client.apiClient.SetHostURL(cfg.ApiUrl)

	client.statusPath = "/users/%s/robots"
	client.insightsPath = client.statusPath + "/%s/insights"
	client.cmdPath = client.statusPath + "/%s/dispatch-commands"

	err := client.login()
	client.apiClient.SetAuthToken(client.token)
	client.apiClient.SetHeader("x-api-key", cfg.ApiKey)
	return client, err
}

// RefreshToken - Refreshes the access_token granted by the initial client creation.
func (c *Client) RefreshToken() {
	loginBody := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {c.ClientId},
		"client_secret": {c.ClientSecret},
		"refresh_token": {c.refreshToken},
	}

	resp, err := c.authClient.R().
		SetFormDataFromValues(loginBody).
		SetResult(loginResponse{}).
		Post("/oauth/token")

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to refresh access token")
	}

	if resp.StatusCode() != http.StatusOK {
		log.WithFields(log.Fields{
			"statusCode": resp.StatusCode(),
			"status":     resp.Status(),
		}).Fatal("Failed to refresh access token")
	}

	lr := resp.Result().(*loginResponse)
	c.apiClient.SetAuthToken(lr.Token)
	c.refreshToken = lr.RefreshToken
	log.Info("Token successfully refreshed.")
}

// States - Fetch states from the Litter Robot API
func (c *Client) States() ([]State, error) {
	results, err := c.getStates()
	if err != nil {
		return nil, err
	}

	robots := make([]State, 0)
	for _, result := range results {
		s := newState(result)
		robots = append(robots, s)
		c.robots[s.LitterRobotID] = s
	}

	return robots, nil
}

// Insights - return the Litter Robot Insights over the specified number of days.
func (c *Client) Insights(id string, days, timezoneOffset int) ([]Insight, error) {
	if days < 1 {
		days = 1 // litter robot api does not check for zero
	}
	insights := make([]Insight, 0)

	log.WithField("robot", id).Debug("Fetching insights")
	resp, err := c.apiClient.R().
		SetQueryParam("days", fmt.Sprintf("%d", days)).
		SetQueryParam("timezoneOffset", fmt.Sprintf("%d", timezoneOffset)).
		SetResult(Insight{}).
		Get(fmt.Sprintf(c.insightsPath, c.userID, id))

	if err != nil {
		log.WithField("error", err).WithField("robot", id).Error("Failed getting robot insights")
	}

	ins := resp.Result().(*Insight)

	insights = append(insights, *ins)
	log.WithField("robot", id).Info("Successfully fetched insights")
	return insights, nil
}

func (c *Client) login() error {
	log.Debug("Trying to login to the Litter Robot API")

	loginBody := url.Values{
		"username":      {c.Email},
		"password":      {c.Password},
		"grant_type":    {"password"},
		"client_id":     {c.ClientId},
		"client_secret": {c.ClientSecret},
	}

	resp, err := c.authClient.R().
		SetFormDataFromValues(loginBody).
		SetResult(loginResponse{}).
		Post("/oauth/token")

	if err != nil {
		return err
	}

	log.Debug("Login to the Litter Robot API succeeded")
	result := resp.Result().(*loginResponse)

	c.token = result.Token
	type claims struct {
		UserId string `json:"userId"`
		jwt.StandardClaims
	}

	t, err := jwt.ParseWithClaims(result.Token, &claims{}, nil)
	if err != nil {
		log.WithField("error", err.Error()).Error("failed to parse token claims")
		return err
	}
	c.userID = t.Claims.(*claims).UserId
	c.Expiry = time.Second * time.Duration(result.ExpiresIn)
	c.refreshToken = result.RefreshToken

	return nil
}

func (c *Client) getStates() ([]robotResponse, error) {
	log.Debug("Trying to fetch data from the Litter Robot API")

	if c.token == "" {
		if err := c.login(); err != nil {
			return nil, err
		}
	}

	resp, err := c.apiClient.R().
		SetResult([]robotResponse{}).
		Get(fmt.Sprintf(c.statusPath, c.userID))

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to fetch robot information")
		return nil, err
	}

	rresp := resp.Result().(*[]robotResponse)

	log.Debug("Fetching data from the Litter Robot API succeeded")
	return *rresp, nil
}

func (c *Client) sendCommand(robotId string, command string) error {
	// Fetch a token if we don't have one
	if c.token == "" {
		if err := c.login(); err != nil {
			return err
		}
	}

	log.WithFields(log.Fields{
		"command": command,
		"robotId": robotId,
	}).Info("Sending command to Litter Robot")

	body, _ := json.Marshal(map[string]string{
		"command":       command,
		"litterRobotId": robotId,
	})

	result, err := c.apiClient.R().SetBody(body).Post(fmt.Sprintf(c.cmdPath, c.userID, robotId))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error sending command to Litter Robot")
		return err
	}

	if result.StatusCode() != http.StatusOK {
		log.WithFields(log.Fields{
			"statusCode": result.StatusCode(),
			"status":     result.Status(),
			"id":         robotId,
		}).Error("Command request failed")
		return errors.New(result.Status())
	}

	return nil
}

// PowerOn - Turn unit power on.
func (c *Client) PowerOn(robotId string) {
	log.Debug("Sending Power On Command")
	_ = c.sendCommand(robotId, powerCmd+on)
}

// PowerOff - Turn unit power off.
func (c *Client) PowerOff(robotId string) {
	log.Debug("Sending Power Off Command")
	_ = c.sendCommand(robotId, powerCmd+off)
}

// NightLightOn - Turn nightlight on.
func (c *Client) NightLightOn(robotId string) {
	log.Debug("Sending Night Light On Command")
	_ = c.sendCommand(robotId, nightLightCmd+on)
}

// NightLightOff - Turn nightlight off.
func (c *Client) NightLightOff(robotId string) {
	log.Debug("Sending Night Light Off Command")
	_ = c.sendCommand(robotId, nightLightCmd+off)
}

// PanelLockOn - Enable the panel lock.
func (c *Client) PanelLockOn(robotId string) {
	log.Debug("Sending Panel Lock On Command")
	_ = c.sendCommand(robotId, panelLockCmd+on)
}

// PanelLockOff - Disable the panel lock.
func (c *Client) PanelLockOff(robotId string) {
	log.Debug("Sending Panel Lock Off Command")
	_ = c.sendCommand(robotId, panelLockCmd+off)
}

// Cycle - Start a clean cycle.
func (c *Client) Cycle(robotId string) {
	log.Debug("Sending Cycle Command")
	_ = c.sendCommand(robotId, cycleCmd)
}

// Wait - Set clean cycle wait time.
func (c *Client) Wait(robotId string, val string) {
	log.Debug("Sending Wait Command")
	_ = c.sendCommand(robotId, waitCmd+val)
}
