#!/bin/bash

PLATFORMS=("windows" "linux")
ARCHS=("amd64" "arm64" "386")
OUTPUT_DIR="builds"

mkdir -p "$OUTPUT_DIR"

for platform in "${PLATFORMS[@]}"; do
    for arch in "${ARCHS[@]}"; do
        output_name="postfix-to-cloudflare-$platform-$arch"
        if [ "$arch" == "arm64" ] && [ "$platform" == "linux" ]; then
            # Set GOARM to 7 for armv7, but only for Linux
            GOARM="7"
            output_name+="-v7"
        else
            GOARM=""
        fi
        # Append .exe for Windows builds
        if [ "$platform" == "windows" ]; then
            output_name+=".exe"
        fi
        env GOOS=$platform GOARCH=$arch GOARM=$GOARM go build -o "$OUTPUT_DIR/$output_name" ./main
    done
done