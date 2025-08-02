#!/bin/bash

# Backend Deployment Script
echo "🚀 Preparing Backend for Deployment..."

# Check if we're in the Backend directory
if [ ! -f "main.go" ]; then
    echo "❌ Error: Please run this script from the Backend directory"
    exit 1
fi

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed or not in PATH"
    exit 1
fi

# Test build
echo "🔨 Testing build..."
go build -o main .

if [ -f "main" ] || [ -f "main.exe" ]; then
    echo "✅ Build successful! Backend is ready for deployment."
    rm -f main main.exe  # Clean up test build
    echo ""
    echo "📋 Next steps:"
    echo "1. Push your code to GitHub"
    echo "2. Import the repository in Vercel"
    echo "3. Set environment variables:"
    echo "   - GO_VERSION=1.21"
    echo "   - GIN_MODE=release"
    echo "   - ALLOWED_ORIGINS=*"
    echo "4. Deploy!"
else
    echo "❌ Build failed! Please check for errors."
    exit 1
fi
