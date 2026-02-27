#!/usr/bin/env bash
# Minimal installer for Ubuntu (assumes user runs as root or with sudo)
set -euo pipefail

echo "Updating apt and installing prerequisites..."
apt-get update
apt-get install -y build-essential git curl ca-certificates

GO_VERSION=${GO_VERSION:-1.22.0}
echo "Checking Go installation..."
if command -v go >/dev/null 2>&1; then
  echo "go found: $(go version)"
  GO_BIN=$(command -v go)
else
  echo "Go not found; installing Go ${GO_VERSION}"
  archive="go${GO_VERSION}.linux-amd64.tar.gz"
  url="https://go.dev/dl/${archive}"
  echo "Downloading ${url} ..."
  if ! curl -fSL "$url" -o /tmp/go.tgz; then
    echo "failed to download ${url}."
    echo "Please install Go manually (apt, snap, or from https://go.dev/dl/) or set GO_VERSION to a valid release and re-run the installer."
    exit 1
  fi
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf /tmp/go.tgz
  export PATH=$PATH:/usr/local/go/bin
  GO_BIN=/usr/local/go/bin/go
fi

echo "Building smtp-ingress-service..."
cd /opt/viramail || mkdir -p /opt/viramail && cd /opt/viramail
# Clone the provided repository (will reuse if already present)
if [ -d .git ]; then
  echo "Existing repo found, fetching..."
  git fetch --all
  git reset --hard origin/main || true
else
  git clone https://github.com/viramail-viragar/viramail.git .
fi

# configure git user/email for commits (installer-level)
git config user.email "viramail@viragar.ir" || true
git config user.name "ViraMail Installer" || true

/usr/bin/env ${GO_BIN:-go} build -o /usr/local/bin/viramail-smtp ./cmd/smtp-ingress-service

echo "Create systemd service..."
cat >/etc/systemd/system/viramail-smtp.service <<'EOS'
[Unit]
Description=ViraMail SMTP Ingress Service
After=network.target

[Service]
Environment="CERT_PATH=/etc/letsencrypt/live/viragar.ir/fullchain.pem"
Environment="KEY_PATH=/etc/letsencrypt/live/viragar.ir/privkey.pem"
ExecStart=/usr/local/bin/viramail-smtp
Restart=on-failure
User=www-data

[Install]
WantedBy=multi-user.target
EOS

systemctl daemon-reload
systemctl enable --now viramail-smtp.service
echo "Installation complete."

echo "Attempting to obtain Let's Encrypt certificate for viragar.ir using certbot (standalone)."
if ! command -v certbot >/dev/null 2>&1; then
  apt-get update
  apt-get install -y snapd
  snap install core; snap refresh core
  snap install --classic certbot
  ln -s /snap/bin/certbot /usr/bin/certbot || true
fi

if [ -f /etc/letsencrypt/live/viragar.ir/fullchain.pem ]; then
  echo "Certificate already exists at /etc/letsencrypt/live/viragar.ir/"
else
  echo "Stopping service to allow certbot standalone to bind to port 80/443..."
  systemctl stop viramail-smtp.service || true
  certbot certonly --standalone --non-interactive --agree-tos -m viramail@viragar.ir -d viragar.ir || true
  echo "Starting service"
  systemctl start viramail-smtp.service || true
fi
