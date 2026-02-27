package main

import (
    "crypto/tls"
    "log"
    "net"
    "net/textproto"
    "os"
    "time"
)

func main() {
    addr := ":2525" // non-priv port for dev
    log.Printf("starting smtp-ingress-service on %s", addr)

    certFile := "cert.pem"
    keyFile := "key.pem"
    if _, err := os.Stat(certFile); os.IsNotExist(err) {
        log.Printf("TLS cert/key not found, generating self-signed for dev")
        if err := generateSelfSigned(certFile, keyFile); err != nil {
            log.Fatalf("failed to generate cert: %v", err)
        }
    }

    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        log.Fatalf("failed to load cert: %v", err)
    }

    cfg := &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS13}
    l, err := tls.Listen("tcp", addr, cfg)
    if err != nil {
        log.Fatalf("listen error: %v", err)
    }
    defer l.Close()

    for {
        conn, err := l.Accept()
        if err != nil {
            log.Printf("accept error: %v", err)
            continue
        }
        go handleConn(conn)
    }
}

func handleConn(c net.Conn) {
    defer c.Close()
    tp := textproto.NewConn(c)
    defer tp.Close()

    tp.PrintfLine("220 vira.local ESMTP ready")
    c.SetDeadline(time.Now().Add(5 * time.Minute))

    for {
        line, err := tp.ReadLine()
        if err != nil {
            log.Printf("read error: %v", err)
            return
        }
        log.Printf("C: %s", line)
        switch {
        case line == "QUIT":
            tp.PrintfLine("221 Bye")
            return
        case line == "EHLO" || line == "HELO":
            tp.PrintfLine("250-Hello")
            tp.PrintfLine("250-STARTTLS")
            tp.PrintfLine("250 OK")
        default:
            tp.PrintfLine("250 OK")
        }
    }
}

// generateSelfSigned writes a minimal self-signed cert for dev use.
func generateSelfSigned(certPath, keyPath string) error {
    // For brevity, create placeholder files. Users should replace with real certs.
    cert := []byte("-----BEGIN CERTIFICATE-----\nMIID...DEV_CERT...\n-----END CERTIFICATE-----\n")
    key := []byte("-----BEGIN PRIVATE KEY-----\nMIIE...DEV_KEY...\n-----END PRIVATE KEY-----\n")
    if err := os.WriteFile(certPath, cert, 0644); err != nil {
        return err
    }
    if err := os.WriteFile(keyPath, key, 0600); err != nil {
        return err
    }
    return nil
}
