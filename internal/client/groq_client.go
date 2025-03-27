package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type GroqClient struct {
    APIKey  string
    BaseURL string
}

func NewGroqClient(apiKey, baseURL string) *GroqClient {
    return &GroqClient{APIKey: apiKey, BaseURL: baseURL}
}

func (gc *GroqClient) StreamResponse(ctx context.Context, prompt string) (io.ReadCloser, error) {
    payload := map[string]interface{}{
        "model": "llama3-8b-8192",
        "messages": []map[string]string{
            {"role": "user", "content": prompt},
        },
        "stream": true,
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequestWithContext(ctx, "POST", gc.BaseURL+"/v1/chat/completions", bytes.NewReader(body))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Bearer "+gc.APIKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }

    return resp.Body, nil
}