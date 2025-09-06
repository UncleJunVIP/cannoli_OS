#!/bin/bash

set -e

echo "Building RetroArch Docker image..."
docker build -t retroarch-builder .

echo "Creating temporary container..."
CONTAINER_ID=$(docker create retroarch-builder)

rm -rf ./build/RetroArch || true
mkdir -p ./build/RetroArch

echo "Copying files from container..."
docker cp $CONTAINER_ID:/home/builder/out/. ./build/RetroArch/

echo "Cleaning up container..."
docker rm $CONTAINER_ID

echo "Build completed! Files copied to ./build/RetroArch"
ls -la ./build/RetroArch

cp retroarch.cfg ./build/RetroArch
cp launch.sh ./build/RetroArch
