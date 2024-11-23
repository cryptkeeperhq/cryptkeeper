#!/bin/bash

# Exit script on error
set -e

# Directory for certificates
CERT_DIR="./certs"
mkdir -p "$CERT_DIR"

# Generate CA key and certificate
openssl genrsa -out "$CERT_DIR/ca.key" 2048
openssl req -x509 -new -nodes -key "$CERT_DIR/ca.key" -sha256 -days 365 -out "$CERT_DIR/ca.pem" -subj "/C=US/ST=CA/L=SanFrancisco/O=MyOrg/OU=MyCA"

echo "CA certificate and key generated."

# Generate server key and CSR
openssl genrsa -out "$CERT_DIR/server.key" 2048
openssl req -new -key "$CERT_DIR/server.key" -out "$CERT_DIR/server.csr" -subj "/C=US/ST=CA/L=SanFrancisco/O=MyOrg/OU=MyServer"

# Generate server certificate signed by the CA
openssl x509 -req -in "$CERT_DIR/server.csr" -CA "$CERT_DIR/ca.pem" -CAkey "$CERT_DIR/ca.key" -CAcreateserial -out "$CERT_DIR/server.crt" -days 365 -sha256

echo "Server certificate and key generated."

# Cleanup
rm "$CERT_DIR/server.csr"

echo "Certificates are in the $CERT_DIR directory."
