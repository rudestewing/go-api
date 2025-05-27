#!/bin/bash

# Development setup script for Go API

echo "🚀 Setting up Go API development environment..."

# Check if .env exists
if [ ! -f .env ]; then
    echo "📄 Creating .env file from template..."
    cp .env.example .env
    echo "✅ .env file created. Please edit it with your configuration."
    echo "⚠️  Don't forget to set DATABASE_URL and JWT_SECRET!"
else
    echo "✅ .env file already exists"
fi

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy

# Build the application
echo "🔨 Building application..."
go build -o tmp/main .

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "🎉 Setup complete! To start the server:"
    echo "   go run main.go"
    echo ""
    echo "📚 Or use Air for hot reload:"
    echo "   air"
else
    echo "❌ Build failed. Please check your configuration."
    exit 1
fi
