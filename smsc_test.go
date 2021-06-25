package smsc_test

import (
	"bytes"
	"errors"
	"github.com/larship/smsc"
	"io/ioutil"
	"net/http"
	"testing"
)

type TestHTTPClient struct {
	errorMode bool
}

func (c *TestHTTPClient) SetErrorMode(state bool) {
	c.errorMode = state
}

func (c *TestHTTPClient) Do(req *http.Request) (*http.Response, error) {
	json := `{}`
	nopCloser := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	if c.errorMode {
		return nil, errors.New("Error")
	}

	return &http.Response{
		StatusCode: 200,
		Body:       nopCloser,
	}, nil
}

func TestClient_SendSms(t *testing.T) {
	testHttpClient := &TestHTTPClient{}

	client := smsc.New("test", "test")
	client.SetHTTPClient(testHttpClient)
	resp, err := client.SendSms("test", "test")

	if resp == nil || err != nil {
		t.Fatal("Error occurred when testing request without errors")
	}

	testHttpClient.SetErrorMode(true)
	resp, err = client.SendSms("test", "test")

	if resp != nil || err == nil {
		t.Fatal("Error occurred when testing request with errors")
	}
}
