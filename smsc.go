package smsc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const smscUrl = "https://smsc.ru/sys/send.php"

type Client struct {
	client       HttpClient
	login        string
	password     string
	sender       string
	sender_email string
	tinyurl      string
	charset      string
	format       string
}

type Response struct {
	Id        int    `json:"id"`
	Count     int    `json:"cnt"`
	ErrorText string `json:"error"`
	ErrorCode int    `json:"error_code"`
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Create new API client instance
func New(login string, password string) *Client {
	client := &Client{
		client:   &http.Client{},
		login:    login,
		password: password,
		charset:  "utf-8",
		format:   "3",
	}

	return client
}

// Set HTTP client
func (c *Client) SetHTTPClient(http HttpClient) {
	c.client = http
}

// Set API login
func (c *Client) SetLogin(login string) {
	c.login = login
}

// Set API password
func (c *Client) SetPassword(password string) {
	c.password = password
}

// Set message sender parameter
func (c *Client) SetSender(sender string) {
	c.sender = sender
}

// Set message charset parameter
func (c *Client) SetCharset(charset string) {
	c.charset = charset
}

// Set message format parameter
func (c *Client) SetFormat(format string) {
	c.format = format
}

// Send message
func (c *Client) send(params *url.Values) (*Response, error) {
	req, _ := http.NewRequest("POST", smscUrl, bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var respJson Response
	err = json.NewDecoder(resp.Body).Decode(&respJson)

	if err != nil {
		return nil, err
	}

	if respJson.ErrorCode != 0 || respJson.ErrorText != "" {
		return nil, fmt.Errorf("Send SMS error: %d, %s", respJson.ErrorCode, respJson.ErrorText)
	}

	return &respJson, nil
}

// Send SMS
func (c *Client) SendSms(phone string, text string) (*Response, error) {
	params := url.Values{}
	params.Add("login", c.login)
	params.Add("psw", c.password)
	params.Add("phones", phone)
	params.Add("mes", text)
	params.Add("sender", c.sender)
	params.Add("charset", c.charset)
	params.Add("fmt", c.format)

	return c.send(&params)
}
