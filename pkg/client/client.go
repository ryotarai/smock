package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	EventURL      string
	ExternalURL   string
	UserName      string
	UserID        string
	SigningSecret string
}

func New() *Client {
	return &Client{}
}

func (c *Client) SendCommand(command, text string) (*slack.Msg, error) {
	responseURL := c.ExternalURL
	if !strings.HasSuffix(responseURL, "/") {
		responseURL += "/"
	}
	responseURL += "a/response"

	form := url.Values{}
	form.Set("command", command)
	form.Set("text", text)
	form.Set("response_url", responseURL)
	form.Set("user_id", c.UserID)
	form.Set("user_name", c.UserName)
	form.Set("team_id", "TEAMID")
	form.Set("channel_id", "CHANNELID")
	body := form.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", c.EventURL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	err = c.signRequest(req)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return nil, fmt.Errorf("invalid response: %d", resp.StatusCode)
	}

	parts := strings.SplitN(resp.Header.Get("content-type"), ";", 2)
	contentType := parts[0]

	msg := &slack.Msg{}

	switch contentType {
	case "text/html":
	case "application/json":
		json.NewDecoder(resp.Body).Decode(msg)
	default:
		return nil, fmt.Errorf("unexpected content-type: %s", contentType)
	}

	return msg, nil
}

type messageEvent struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	TS      string `json:"ts"`
}

type event struct {
	Token       string      `json:"token"`
	TeamID      string      `json:"team_id"`
	ApiAppID    string      `json:"api_app_id"`
	AuthedUsers []string    `json:"authed_users"`
	Type        string      `json:"type"`
	Event       interface{} `json:"event"`
	EventID     string      `json:"event_id"`
	EventTime   int         `json:"event_time"`
}

func (c *Client) SendMessage(text string) error {
	e := event{
		Token:       "TOKEN",
		TeamID:      "TEAMID",
		ApiAppID:    "APIAPPID",
		AuthedUsers: []string{},
		Type:        "event_callback",
		Event: messageEvent{
			Type:    "message",
			Text:    text,
			Channel: "CHANNELID",
			User:    c.UserID,
			TS:      strconv.FormatInt(time.Now().UnixNano(), 10),
		},
		EventID:   strconv.FormatInt(time.Now().UnixNano(), 10),
		EventTime: int(time.Now().Unix()),
	}

	b, err := json.Marshal(&e)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.EventURL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	err = c.signRequest(req)
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return fmt.Errorf("invalid response: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) signRequest(req *http.Request) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	ts := time.Now().Unix()
	signature, err := c.calcSignature(ts, string(body))
	if err != nil {
		return err
	}

	req.Header.Set("X-Slack-Request-Timestamp", strconv.FormatInt(ts, 10))
	req.Header.Set("X-Slack-Signature", signature)
	return nil
}

func (c *Client) calcSignature(ts int64, requestBody string) (string, error) {
	base := fmt.Sprintf("v0:%d:%s", ts, requestBody)
	mac := hmac.New(sha256.New, []byte(c.SigningSecret))
	_, err := fmt.Fprint(mac, base)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("v0=%x", mac.Sum(nil)), nil
}
