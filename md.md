# Navigate to infrastructure
cd C:\devoptics\infrastructure

# Start Docker services
docker-compose up -d

# Wait for services to be ready
Start-Sleep -Seconds 15

# Verify all services are running
docker ps
```

**Expected Output:**
```
✅ devoptics_postgres   - Port 5432
✅ devoptics_redis      - Port 6379
✅ devoptics_mongodb    - Port 27017
# Test PostgreSQL
Test-NetConnection localhost -Port 5432

# Test Redis
Test-NetConnection localhost -Port 6379

# Test MongoDB
Test-NetConnection localhost -Port 27017
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


cd C:\devoptics\infrastructure\database\migrations

migrate -path . -database "postgresql://devoptics:devoptics_pass@localhost:5432/devoptics?sslmode=disable" up

# Verify migrations
migrate -path . -database "postgresql://devoptics:devoptics_pass@localhost:5432/devoptics?sslmode=disable" version
```

**Expected Output:**
```
✅ 3/u init_schema (2026-02-13...)
✅ Database schema created successfully
# Navigate to auth service
cd C:\devoptics\services\auth-service

# Initialize Go module
go mod tidy

# Run the service
go run cmd/server/main.go
```

**Expected Output:**
```
✅ Connected to PostgreSQL
🚀 Auth Service starting on port 8091
📝 Environment: development
[GIN-debug] Listening and serving HTTP on :8091



# Test health check
Invoke-WebRequest http://localhost:8091/health

# Test register
$registerBody = @{
    email = "test@devoptics.com"
    password = "password123"
    first_name = "Test"
    last_name = "User"
} | ConvertTo-Json

Invoke-WebRequest -Uri http://localhost:8091/api/v1/auth/register `
    -Method POST `
    -ContentType "application/json" `
    -Body $registerBody

# Test login
$loginBody = @{
    email = "test@devoptics.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-WebRequest -Uri http://localhost:8091/api/v1/auth/login `
    -Method POST `
    -ContentType "application/json" `
    -Body $loginBody

# Extract token
$token = ($response.Content | ConvertFrom-Json).token
Write-Host "Token: $token" -ForegroundColor Green



cd C:\devoptics\frontend

# Create Next.js app
npx create-next-app@latest web --typescript --tailwind --app --src-dir --import-alias "@/*"

cd web

# Install dependencies
npm install axios zustand clsx tailwind-merge lucide-react recharts

