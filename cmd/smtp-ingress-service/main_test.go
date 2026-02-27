package main

import (
    "bufio"
    "crypto/tls"
    "strings"
    "testing"
    "time"

    "github.com/viramail/viramail/internal/storage"
)

// This is a lightweight integration test that starts the server in a goroutine
// and issues a minimal SMTP conversation including DATA, then verifies storage.
func TestSMTP_DataSave(t *testing.T) {
    go func() {
        main()
    }()
    // give server a moment to start
    time.Sleep(200 * time.Millisecond)

    // connect using TLS because server listens with TLS
    importTls := true
    _ = importTls
    // create a TLS connection that skips verification for self-signed cert
    // (only for test)
    conn, err := tls.Dial("tcp", "127.0.0.1:2525", &tls.Config{InsecureSkipVerify: true})
    if err != nil {
        t.Fatalf("dial failed: %v", err)
    }
    defer conn.Close()
    r := bufio.NewReader(conn)
    // read banner
    _, _ = r.ReadString('\n')

    send := func(cmd string) string {
        conn.Write([]byte(cmd + "\r\n"))
        s, _ := r.ReadString('\n')
        return s
    }

    send("EHLO test")
    send("MAIL FROM:<a@b.com>")
    send("RCPT TO:<c@d.com>")
    send("DATA")
    // read 354
    // send body
    conn.Write([]byte("Subject: hi\r\n\r\nhello\r\n.\r\n"))
    // read 250
    s, _ := r.ReadString('\n')
    if !strings.Contains(s, "250") {
        t.Fatalf("expected 250 got: %s", s)
    }

    // verify storage
    if storage.LastSavedID == "" {
        t.Fatalf("expected saved id set")
    }
    if !strings.Contains(string(storage.LastSavedData), "hello") {
        t.Fatalf("saved data didn't contain message body")
    }
}
