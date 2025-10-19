#!/bin/bash
# Copyright 2025 Paolo Reyes
# Licensed under GNU Affero General Public License v3.0 (AGPL-3.0)
#
# SSL Certificate Renewal Script for Void Chronicles
# This script renews Let's Encrypt certificates and copies them to the Docker volume

set -e

DOMAIN="vc.internetblacksmith.dev"
CONTAINER_NAME="void-chronicles-web"
VOLUME_PATH="/var/lib/docker/volumes/void-ssl/_data"

echo "==> Stopping Void Chronicles container for certificate renewal..."
docker stop $(docker ps -q --filter "name=${CONTAINER_NAME}") || echo "No container running"

echo "==> Renewing Let's Encrypt certificate for ${DOMAIN}..."
certbot renew --standalone --non-interactive

echo "==> Copying renewed certificates to Docker volume..."
sudo cp /etc/letsencrypt/live/${DOMAIN}/fullchain.pem ${VOLUME_PATH}/cert.pem
sudo cp /etc/letsencrypt/live/${DOMAIN}/privkey.pem ${VOLUME_PATH}/key.pem
sudo chmod 644 ${VOLUME_PATH}/cert.pem
sudo chmod 644 ${VOLUME_PATH}/key.pem

echo "==> Starting Void Chronicles container..."
docker start $(docker ps -aq --filter "name=${CONTAINER_NAME}" | head -1)

echo "==> Certificate renewal complete!"
echo "==> Certificate valid until: $(openssl x509 -enddate -noout -in ${VOLUME_PATH}/cert.pem)"
