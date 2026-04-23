#!/bin/bash

# SSL Certificate Generation Script for Development
# This script generates self-signed certificates for local development

set -e

CERT_DIR="./certs"
CERT_FILE="$CERT_DIR/server.crt"
KEY_FILE="$CERT_DIR/server.key"
DAYS_VALID=365

echo "=== SSL Certificate Generator ==="
echo ""

# Create certs directory if it doesn't exist
if [ ! -d "$CERT_DIR" ]; then
    echo "Creating certificates directory: $CERT_DIR"
    mkdir -p "$CERT_DIR"
fi

# Check if certificates already exist
if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ]; then
    echo "Certificates already exist:"
    echo "  Certificate: $CERT_FILE"
    echo "  Key: $KEY_FILE"
    echo ""
    read -p "Do you want to regenerate them? (y/N) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Keeping existing certificates."
        exit 0
    fi
    echo "Removing old certificates..."
    rm -f "$CERT_FILE" "$KEY_FILE"
fi

echo "Generating self-signed SSL certificate..."
echo ""

# Generate private key and certificate
openssl req -x509 -newkey rsa:4096 -keyout "$KEY_FILE" -out "$CERT_FILE" -days $DAYS_VALID -nodes \
    -subj "/C=CN/ST=Shenzhen/L=Shenzhen/O=AI Parking/OU=Development/CN=localhost" \
    -addext "subjectAltName=DNS:localhost,DNS:*.localhost,IP:127.0.0.1,IP:::1"

echo ""
echo "✅ Certificates generated successfully!"
echo ""
echo "Certificate details:"
echo "  File: $CERT_FILE"
echo "  Key: $KEY_FILE"
echo "  Valid for: $DAYS_VALID days"
echo ""
echo "Certificate info:"
openssl x509 -in "$CERT_FILE" -noout -text | grep -A 2 "Validity"
echo ""
echo "⚠️  IMPORTANT: These are self-signed certificates for development only."
echo "   Browsers will show security warnings - this is normal for self-signed certs."
echo ""
echo "To trust this certificate in your system:"
echo "  macOS: sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain $CERT_FILE"
echo "  Linux:  sudo cp $CERT_FILE /usr/local/share/ca-certificates/ && sudo update-ca-certificates"
echo ""
echo "Server will now run with HTTPS enabled."
