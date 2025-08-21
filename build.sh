#!/bin/bash

# Build script for Fuzzy web server
# This script builds the fuzzy binary and verifies it was created successfully

set -e

echo "Building Fuzzy web server..."

# Build the binary
go build -o fuzzy main.go

# Verify the binary was created
if [ -f "./fuzzy" ]; then
    echo "✓ Build successful! Binary created: ./fuzzy"
    echo "File info:"
    ls -la fuzzy
    echo ""
    echo "To run the server:"
    echo "  ./fuzzy"
    echo ""
    echo "The server will start on http://localhost:8080"
else
    echo "✗ Build failed - binary not found"
    exit 1
fi