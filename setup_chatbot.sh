#!/bin/bash

# ==== CONFIG ====
APP_NAME="chatbot-server"
SERVICE_NAME="chatbot"
INSTALL_DIR="/opt/chatbot"
ENV_FILE=".env"
# ================

echo "📦 Building Go binary..."
if [ ! -f "main.go" ]; then
  echo "❌ main.go tidak ditemukan di direktori ini. Jalankan dari root project."
  exit 1
fi

go mod tidy
go build -o $APP_NAME main.go
if [ $? -ne 0 ]; then
  echo "❌ Build gagal. Cek error di atas."
  exit 1
fi

echo "📁 Deploy binary & env to $INSTALL_DIR..."
sudo mkdir -p $INSTALL_DIR
sudo mv -f $APP_NAME $INSTALL_DIR/
sudo cp -f $ENV_FILE $INSTALL_DIR/
sudo chmod +x $INSTALL_DIR/$APP_NAME

echo "📝 Membuat systemd service..."

sudo tee /etc/systemd/system/${SERVICE_NAME}.service > /dev/null <<EOF
[Unit]
Description=Go Chatbot Service
After=network.target

[Service]
Type=simple
ExecStart=$INSTALL_DIR/$APP_NAME
WorkingDirectory=$INSTALL_DIR
EnvironmentFile=$INSTALL_DIR/$ENV_FILE
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

echo "🔁 Restart systemd service..."
sudo systemctl daemon-reexec
sudo systemctl daemon-reload
sudo systemctl enable $SERVICE_NAME
sudo systemctl restart $SERVICE_NAME

echo "✅ Selesai!"
echo "👉 Cek status: sudo systemctl status $SERVICE_NAME"
echo "👉 Lihat log:  sudo journalctl -u $SERVICE_NAME -f"