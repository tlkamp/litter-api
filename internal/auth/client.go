package auth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

const (
	authURL = "https://42nk7qrhdg.execute-api.us-east-1.amazonaws.com/prod/login"
	apiKey  = "w2tPFbjlP13GUmb8dMjUL5B2YyPVD3pJ7Ey6fz8v"

	oauthURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken"
	oauthKey = "QUl6YVN5Q3Y4NGplbDdKa0NRbHNncXJfc2xYZjNmM3gtY01HMTVR"

	iosHeader      = "x-ios-bundle-identifier"
	iosHeaderValue = "com.whisker.ios"

	refreshEndpoint = "https://securetoken.googleapis.com/v1/token"
)

var iosHeaders = map[string][]string{
	iosHeader: []string{iosHeaderValue},
}

// claims represents the relevant claims contained in the authResponse token.
type claims struct {
	Subclaim struct {
		UserId string `json:"mid"`
	} `json:"claims"`
	jwt.StandardClaims
}

// authResponse is the initial authentication response
type authResponse struct {
	Token string `json:"token"`
}

// authInfo represents the data necessary to log into the LitterRobot API.
type authInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// oauthBody is the payload structure for the oauth token exchange flow.
type oauthBody struct {
	ReturnSecureToken bool   `json:"returnSecureToken"`
	Token             string `json:"token"`
}

// oauthResponse is the response structure for the ouath token exchange flow.
type oauthResponse struct {
	Kind         string `json:"kind"`
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	IsNewUser    bool   `json:"isNewUser"`
}

type refreshBody struct {
	GrantType    string `json:"grantType"`
	RefreshToken string `json:"refreshToken"`
}

type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

// Client interacts with the authN and authZ systems behind the LitterRobot API
type Client struct {
	mutex    *sync.RWMutex
	authInfo *authInfo
	userId   string
	idToken  string
	refToken string
}

// New reutrns an initialized Client.
func New(email, password string) *Client {
	return &Client{
		mutex: &sync.RWMutex{},
		authInfo: &authInfo{
			Email:    email,
			Password: password,
		},
	}
}

// UserID returns the logged in user's User ID.
func (c *Client) UserID() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.userId
}

// IDToken returns the ID token obtained from the Oauth flow.
func (c *Client) IDToken() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.idToken
}

// RefreshToken returns the refresh token obtained from the Oauth flow.
func (c *Client) refreshToken() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.refToken
}

// Login performs the initial login and exchanges the authN token for an access token.
func (c *Client) Login(ctx context.Context) error {
	authToken, err := c.getInitialToken(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting initial token")
	}

	t, err := jwt.ParseWithClaims(authToken, &claims{}, nil)
	if err != nil && !strings.Contains(err.Error(), "Keyfunc") {
		return errors.Wrap(err, "error parsing jwt")
	}

	c.userId = t.Claims.(*claims).Subclaim.UserId

	if err := c.doOauthFlow(ctx, authToken); err != nil {
		return errors.Wrap(err, "error doing oauth flow")
	}

	return nil
}

func (c *Client) getInitialToken(ctx context.Context) (string, error) {
	authInfo, err := json.Marshal(c.authInfo)
	if err != nil {
		return "", errors.Wrap(err, "marshaling auth information")
	}

	auth := bytes.NewReader(authInfo)

	header := map[string][]string{
		"x-api-key": []string{apiKey},
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, authURL, auth)
	req.Header = header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "error calling login")
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "error reading body")
	}

	var authResp *authResponse
	err = json.Unmarshal(respBody, &authResp)
	if err != nil {
		return "", errors.Wrap(err, "unmarshalling auth response")
	}

	return authResp.Token, nil
}

func (c *Client) doOauthFlow(ctx context.Context, authToken string) error {
	decode, err := base64.StdEncoding.DecodeString(oauthKey)
	if err != nil {
		return errors.Wrap(err, "error decoding key")
	}

	decoded := string(decode)

	params := url.Values{}
	params.Set("key", decoded)

	oauthBody := oauthBody{
		ReturnSecureToken: true,
		Token:             authToken,
	}

	data, err := json.Marshal(oauthBody)
	if err != nil {
		return errors.Wrap(err, "error marshaling oauth request body")
	}

	body := bytes.NewBuffer(data)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, oauthURL+"?"+params.Encode(), body)
	if err != nil {
		return errors.Wrap(err, "error creating oauth request")
	}

	req.Header = iosHeaders

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error sending oauth request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected oauth response code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading response body")
	}

	var oauthResp *oauthResponse
	if err := json.Unmarshal(b, &oauthResp); err != nil {
		return errors.Wrap(err, "error unmarshaling oauth response")
	}

	c.mutex.Lock()
	c.idToken = oauthResp.IDToken
	c.refToken = oauthResp.RefreshToken
	c.mutex.Unlock()

	return nil
}

func (c *Client) DoRefreshToken(ctx context.Context) error {
	decode, err := base64.StdEncoding.DecodeString(oauthKey)
	if err != nil {
		return errors.Wrap(err, "error decoding key")
	}

	decoded := string(decode)

	params := url.Values{}
	params.Set("key", decoded)

	rr := &refreshBody{
		GrantType:    "refresh_token",
		RefreshToken: c.refreshToken(),
	}

	body, err := json.Marshal(rr)
	if err != nil {
		return errors.Wrap(err, "error marshaling refresh body")
	}

	payload := bytes.NewBuffer(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, refreshEndpoint+"?"+params.Encode(), payload)
	if err != nil {
		return errors.Wrap(err, "error creating refresh request")
	}

	req.Header = iosHeaders

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error refreshing token")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected oauth response code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading response body")
	}

	var refResp *refreshResponse
	if err = json.Unmarshal(b, &refResp); err != nil {
		return errors.Wrap(err, "error unmarshaling refresh response")
	}

	c.mutex.Lock()
	c.idToken = refResp.IDToken
	c.mutex.Unlock()

	return nil
}
