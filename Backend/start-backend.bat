@echo off
echo Checking if Go is installed...
go version
if %errorlevel% equ 0 (
    echo Go is installed successfully!
    echo.
    echo Running backend setup commands...
    echo.
    
    echo 1. Installing dependencies...
    go mod tidy
    
    echo.
    echo 2. Seeding sample data...
    go run seed.go
    
    echo.
    echo 3. Starting the backend server...
    echo Server will run on http://localhost:8080
    echo Press Ctrl+C to stop the server
    echo.
    go run main.go
) else (
    echo Go is not installed or not in PATH.
    echo Please install Go from: https://golang.org/dl/
    echo After installation, run this script again.
    pause
)
