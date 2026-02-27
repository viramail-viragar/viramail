package storage

import (
    "context"
    "errors"
)

type Client struct {
    // placeholder for DB/S3 clients
}

func NewClient() *Client { return &Client{} }

func (c *Client) SaveMail(ctx context.Context, data []byte) (string, error) {
    if len(data) == 0 {
        return "", errors.New("empty mail")
    }
    // stub: return fake id
    return "mail_12345", nil
}
