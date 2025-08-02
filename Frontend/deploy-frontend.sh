#!/bin/bash

# Frontend Deployment Script
echo "🚀 Preparing Frontend for Deployment..."

# Check if we're in the Frontend directory
if [ ! -f "package.json" ]; then
    echo "❌ Error: Please run this script from the Frontend directory"
    exit 1
fi

# Install dependencies
echo "📦 Installing dependencies..."
npm install

# Build the project
echo "🔨 Building project..."
npm run build

# Check if build was successful
if [ -d "dist" ]; then
    echo "✅ Build successful! Frontend is ready for deployment."
    echo ""
    echo "📋 Next steps:"
    echo "1. Push your code to GitHub"
    echo "2. Import the repository in Vercel"
    echo "3. Set environment variable: VITE_API_BASE_URL"
    echo "4. Deploy!"
else
    echo "❌ Build failed! Please check for errors."
    exit 1
fi
