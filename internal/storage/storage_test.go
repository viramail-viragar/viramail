package storage

import (
    "context"
    "testing"
)

func TestSaveMail_Empty(t *testing.T) {
    c := NewClient()
    _, err := c.SaveMail(context.Background(), []byte{})
    if err == nil {
        t.Fatalf("expected error for empty mail")
    }
}

func TestSaveMail_Success(t *testing.T) {
    c := NewClient()
    id, err := c.SaveMail(context.Background(), []byte("hello"))
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if id == "" {
        t.Fatalf("expected non-empty id")
    }
}
