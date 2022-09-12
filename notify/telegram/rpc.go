package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/eoozs/befirst/pkg/net"
)

type APIWrapper interface {
	SendMessage(chatID int64, msg string) error
	GetUpdates() ([]Update, error)
}

type RPCApiWrapper struct {
	BaseURL    string
	BotToken   string
	HTTPClient net.HTTPClient
}

func NewRPCApiWrapper(baseURL, botToken string, httpClient net.HTTPClient) *RPCApiWrapper {
	return &RPCApiWrapper{
		BaseURL:    baseURL,
		BotToken:   botToken,
		HTTPClient: httpClient,
	}
}

const (
	rpcMethodSendMessage = "sendMessage"
	rpcMethodGetUpdates  = "getUpdates"
)

func (c RPCApiWrapper) SendMessage(chatID int64, msg string) error {
	reqBody, _ := json.Marshal(ReqPayloadSendMessage{
		ChatID:                chatID,
		Text:                  msg,
		DisableWebPagePreview: true,
	})

	resp, err := c.HTTPClient.Post(
		c.buildURLOfFunc(rpcMethodSendMessage),
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return fmt.Errorf("unable to do POST request: %v", err)
	}

	if _, err = parseResponse[Message](*resp); err != nil {
		return fmt.Errorf("unable to parse POST response: %v", err)
	}

	return nil
}

func (c RPCApiWrapper) GetUpdates() ([]Update, error) {
	resp, err := c.HTTPClient.Get(
		c.buildURLOfFunc(rpcMethodGetUpdates),
	)
	if err != nil {
		return []Update{}, fmt.Errorf("unable to do GET request: %v", err)
	}

	response, err := parseResponse[[]Update](*resp)
	if err != nil {
		return []Update{}, fmt.Errorf("unable to parse GET response: %v", err)
	}

	return *response, nil
}

func (c RPCApiWrapper) buildURLOfFunc(functionName string) string {
	return fmt.Sprintf("%s/bot%s/%s", c.BaseURL, c.BotToken, functionName)
}

func parseResponse[T any](resp http.Response) (*T, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("request failed")
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response reading failed: %v", err)
	}

	genericResp := RPCGenericResponse{}
	err = json.Unmarshal(respBytes, &genericResp)
	if err != nil {
		return nil, fmt.Errorf("response parsing failed: %v", err)
	}

	if !genericResp.OK {
		return nil, errors.New("response contains invalid data")
	}

	result := new(T)
	err = json.Unmarshal(genericResp.Result, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to parse result: %v", err)
	}

	return result, nil
}
