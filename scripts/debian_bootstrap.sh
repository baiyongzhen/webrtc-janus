#!/bin/bash

set -e

# Install Docker
curl -fsSL get.docker.com -o get-docker.sh
sh get-docker.sh
sudo usermod -aG docker vagrant

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/1.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install Gstreamer
sudo apt-get install -y libgstreamer1.0-dev \
    libgstreamer-plugins-base1.0-dev \
    libgstreamer-plugins-bad1.0-dev \
    gstreamer1.0-plugins-base \
    gstreamer1.0-plugins-good \
    gstreamer1.0-plugins-bad \
    gstreamer1.0-plugins-ugly \
    gstreamer1.0-libav \
    gstreamer1.0-doc \
    gstreamer1.0-tools \
    gstreamer1.0-x \
    gstreamer1.0-alsa \
    gstreamer1.0-gl \
    gstreamer1.0-gtk3 \
    gstreamer1.0-qt5 \
    gstreamer1.0-pulseaudio

export PKG_CONFIG_PATH="/usr/local/opt/libffi/lib/pkgconfig"

# Install janus tools
sudo apt-get install -y janus-tools

# Install ffmpeg
sudo apt install -y ffmpeg

# Install Golang
# golang version number
GO_VERSION=1.15
sudo apt-get install -y curl
sudo curl -fsSL "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" | sudo tar Cxz /usr/local

cat >> /home/vagrant/.profile <<EOF
GOPATH=\\$HOME/go
PATH=/usr/local/go/bin:\\$PATH
export GOPATH PATH
EOF

source /home/vagrant/.profile