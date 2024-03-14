#!/usr/bin/env bash

# change this tag when a new release is made
TAG="v1.2"
IS_MACOS=$(uname -s | grep -i Darwin)
IS_LINUX=$(uname -s | grep -i Linux)
ARCH=$(uname -m)
echo "Detected OS: $IS_MACOS $IS_LINUX $ARCH"
if [[ -n "$IS_MACOS" || -n "$IS_LINUX" ]]; then
	echo "Downloading binary for tag $TAG"
	OS="$([[ -n "$IS_MACOS" ]] && printf Darwin || printf Linux)"
	URL="https://github.com/napisani/json-to-mongo-pipeline-go/releases/download/${TAG}/json-to-mongo-pipeline_${OS}_${ARCH}"
  echo "Downloading from: $URL"
	curl -L -o json-to-mongo-pipeline "$URL"
	chmod +x json-to-mongo-pipeline
	echo "Moving binary to /usr/local/bin -- you may need to enter your password"
	sudo mv json-to-mongo-pipeline /usr/local/bin/json-to-mongo-pipeline
else
	echo "Unsupported OS"
	exit 1
fi
