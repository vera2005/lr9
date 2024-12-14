# Set the output encoding to UTF-8
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

# Start golangci-lint
Write-Host "Running golangci-lint..." -ForegroundColor Cyan
golangci-lint run ./...

# Check the status of the linter execution
if ($LASTEXITCODE -ne 0) {
    Write-Host "Error during golangci-lint execution." -ForegroundColor Red
    exit $LASTEXITCODE
} else {
    Write-Host "golangci-lint completed successfully." -ForegroundColor Green
}

# Compile the project
Write-Host "Compiling the project..." -ForegroundColor Cyan
go build -o myapp.exe .

# Check the status of the compilation
if ($LASTEXITCODE -ne 0) {
    Write-Host "Error during project compilation." -ForegroundColor Red
    exit $LASTEXITCODE
} else {
    Write-Host "Project successfully compiled to myapp.exe." -ForegroundColor Green
}

Write-Host "Script completed successfully." -ForegroundColor Green
