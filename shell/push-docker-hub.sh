#!/bin/bash

VERSION=$1
REPO="edoaurahman/cek-resi-app"
# usage
if [ -z "$VERSION" ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

docker buildx build --platform linux/arm64,linux/amd64  --tag edoaurahman/cek-resi-app:latest . --push
docker buildx build --platform linux/arm64,linux/amd64  --tag edoaurahman/cek-resi-app:$VERSION . --push
