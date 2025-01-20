#!/bin/bash

VERSION=$1
REPO="edoaurahman/cek-resi-app"

docker build -t $REPO:$VERSION -t $REPO:latest .
docker push $REPO:$VERSION
docker push $REPO:latest
