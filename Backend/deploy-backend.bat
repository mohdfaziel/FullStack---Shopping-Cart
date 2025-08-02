@echo off
echo 🚀 Preparing Backend for Deployment...

REM Check if main.go exists
if not exist "main.go" (
    echo ❌ Error: Please run this script from the Backend directory
    pause
    exit /b 1
)

REM Check Go installation
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ Error: Go is not installed or not in PATH
    pause
    exit /b 1
)

REM Test build
echo 🔨 Testing build...
go build -o main.exe .

if exist "main.exe" (
    echo ✅ Build successful! Backend is ready for deployment.
    del main.exe
    echo.
    echo 📋 Next steps:
    echo 1. Push your code to GitHub
    echo 2. Import the repository in Vercel
    echo 3. Set environment variables:
    echo    - GO_VERSION=1.21
    echo    - GIN_MODE=release
    echo    - ALLOWED_ORIGINS=*
    echo 4. Deploy!
) else (
    echo ❌ Build failed! Please check for errors.
    pause
    exit /b 1
)

pause
