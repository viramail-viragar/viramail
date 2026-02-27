package storage

import (
    "context"
    "errors"
)

type Client struct {
    // placeholder for DB/S3 clients
}

func NewClient() *Client { return &Client{} }

// Test hooks (package-level) - allow tests to observe saved data
var LastSavedID string
var LastSavedData []byte

func (c *Client) SaveMail(ctx context.Context, data []byte) (string, error) {
    if len(data) == 0 {
        return "", errors.New("empty mail")
    }
    // store into package-level vars for tests
    LastSavedData = append([]byte(nil), data...)
    LastSavedID = "mail_12345"
    return LastSavedID, nil
}
