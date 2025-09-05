#!/bin/bash
set -e

echo "Building RetroArch Docker image..."
docker build -t retroarch-builder .

mkdir -p output

echo "Running RetroArch build container..."
docker run --rm -v "$(pwd)/output:/home/builder/output" retroarch-builder

echo "Build completed! Checking output directory..."
if [ -f "./output/retroarch" ]; then
    echo "✅ RetroArch binary found!"
    echo "📊 Binary info:"
    file ./output/retroarch
    echo "📏 Binary size:"
    ls -lh ./output/retroarch
    echo ""
    echo "📋 All output files:"
    ls -la ./output/
else
    echo "❌ RetroArch binary not found in output directory"
    echo "Contents of output directory:"
    ls -la ./output/ || echo "Output directory is empty"
    exit 1
fi
