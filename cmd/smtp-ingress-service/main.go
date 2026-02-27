package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/tls"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "log"
    "math/big"
    "net"
    "net/textproto"
    "os"
    "time"
)

func main() {
    addr := ":2525" // non-priv port for dev
    log.Printf("starting smtp-ingress-service on %s", addr)

    // prefer environment variables (set by systemd) then Let's Encrypt default paths
    certFile := os.Getenv("CERT_PATH")
    keyFile := os.Getenv("KEY_PATH")
    if certFile == "" {
        certFile = "/etc/letsencrypt/live/viragar.ir/fullchain.pem"
    }
    if keyFile == "" {
        keyFile = "/etc/letsencrypt/live/viragar.ir/privkey.pem"
    }

    if _, err := os.Stat(certFile); os.IsNotExist(err) {
        log.Printf("TLS cert/key not found at %s, generating self-signed for dev", certFile)
        // create certs locally named cert.pem/key.pem for dev convenience
        certFile = "cert.pem"
        keyFile = "key.pem"
        if err := generateSelfSigned(certFile, keyFile); err != nil {
            log.Fatalf("failed to generate cert: %v", err)
        }
    } else {
        log.Printf("Using TLS cert: %s and key: %s", certFile, keyFile)
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
    // generate a self-signed certificate for development use
    priv, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return err
    }

    serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
    serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
    if err != nil {
        return err
    }

    tmpl := x509.Certificate{
        SerialNumber: serialNumber,
        Subject: pkix.Name{
            Organization: []string{"ViraMail Dev"},
        },
        NotBefore:             time.Now().Add(-1 * time.Hour),
        NotAfter:              time.Now().Add(365 * 24 * time.Hour),
        KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        BasicConstraintsValid: true,
    }

    // include localhost and vira.local
    tmpl.DNSNames = []string{"localhost", "vira.local"}

    derBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
    if err != nil {
        return err
    }

    certOut, err := os.Create(certPath)
    if err != nil {
        return err
    }
    defer certOut.Close()
    if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
        return err
    }

    keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        return err
    }
    defer keyOut.Close()
    privBytes := x509.MarshalPKCS1PrivateKey(priv)
    if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}); err != nil {
        return err
    }
    return nil
}
