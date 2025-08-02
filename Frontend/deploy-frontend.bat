@echo off
echo 🚀 Preparing Frontend for Deployment...

REM Check if package.json exists
if not exist "package.json" (
    echo ❌ Error: Please run this script from the Frontend directory
    pause
    exit /b 1
)

REM Install dependencies
echo 📦 Installing dependencies...
call npm install

REM Build the project
echo 🔨 Building project...
call npm run build

REM Check if build was successful
if exist "dist" (
    echo ✅ Build successful! Frontend is ready for deployment.
    echo.
    echo 📋 Next steps:
    echo 1. Push your code to GitHub
    echo 2. Import the repository in Vercel
    echo 3. Set environment variable: VITE_API_BASE_URL
    echo 4. Deploy!
) else (
    echo ❌ Build failed! Please check for errors.
    pause
    exit /b 1
)

pause
