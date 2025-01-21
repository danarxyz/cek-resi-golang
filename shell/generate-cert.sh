#!/bin/bash

# Check if directory parameter is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <directory>"
    exit 1
fi

DIR=$1

# Check if the directory exists
if [ ! -d "$DIR" ]; then
    echo "Directory $DIR does not exist."
    exit 1
fi

cd "$DIR" || exit

rm *.pem 2>/dev/null

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=ID/ST=East Java/L=Bojonegoro/O=Edodev/OU=Edodev/CN=*.edodev.my.id/emailAddress=edoaurahman@gmail.com"
if [ $? -ne 0 ]; then
    echo "Failed to generate CA's private key and self-signed certificate."
    exit 1
fi

echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=ID/ST=East Java/L=Bojonegoro/O=Edodev/OU=Edodev/CN=*.edodev.my.id/emailAddress=edoaurahman@gmail.com"
if [ $? -ne 0 ]; then
    echo "Failed to generate web server's private key and CSR."
    exit 1
fi

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf
if [ $? -ne 0 ]; then
    echo "Failed to sign web server's CSR."
    exit 1
fi

echo "Server's signed certificate"
openssl x509 -in server-cert.pem -noout -text
