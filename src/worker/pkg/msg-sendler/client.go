package msg_sendler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Client struct {
	BaseURL string
	AuthKey string
}

func NewClient(baseURL string, authKey string) *Client {
	return &Client{BaseURL: baseURL, AuthKey: authKey}
}

type SendMessageRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type SendMessageResponse struct {
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}

func (c *Client) SendMessage(ctx context.Context, content string, phone string) (SendMessageResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", c.BaseURL, nil)
	if err != nil {
		return SendMessageResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("'x-ins-auth-key", c.AuthKey)

	body := SendMessageRequest{
		To:      phone,
		Content: content,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return SendMessageResponse{}, err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	resp, err := client.Do(req)
	if err != nil {
		return SendMessageResponse{}, err
	}

	defer resp.Body.Close()

	var response SendMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return SendMessageResponse{}, err
	}

	return response, nil
}
