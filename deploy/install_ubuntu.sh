#!/usr/bin/env bash
# Minimal installer for Ubuntu (assumes user runs as root or with sudo)
set -euo pipefail

echo "Updating apt and installing prerequisites..."
apt-get update
apt-get install -y build-essential git curl ca-certificates

echo "Installing Go 1.20..."
if ! command -v go >/dev/null 2>&1; then
  curl -fsSL https://go.dev/dl/go1.20.15.linux-amd64.tar.gz -o /tmp/go.tgz
  tar -C /usr/local -xzf /tmp/go.tgz
  export PATH=$PATH:/usr/local/go/bin
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

/usr/local/go/bin/go build -o /usr/local/bin/viramail-smtp ./cmd/smtp-ingress-service

echo "Create systemd service..."
cat >/etc/systemd/system/viramail-smtp.service <<'EOS'
[Unit]
Description=ViraMail SMTP Ingress Service
After=network.target

[Service]
ExecStart=/usr/local/bin/viramail-smtp
Restart=on-failure
User=www-data

[Install]
WantedBy=multi-user.target
EOS

systemctl daemon-reload
systemctl enable --now viramail-smtp.service
echo "Installation complete."
