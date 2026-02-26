# httplog

Simple service to store and retrieve HTTP request data in a PostgreSQL database.

## Usage

```bash
# Start PostgreSQL
docker-compose up -d

# Run the service
go run cmd/main.go

# Insert data
curl -X POST http://localhost:8080/httplog \
  -H "Content-Type: application/json" \
  -d '{
    "id": "a3f5c9d2-7b4e-4c1e-9f8a-2d6b7c8e9f01",
    "url": "/api/v1/orders/12345",
    "method": "GET",
    "time_in": "2026-02-26T14:32:10.123Z",
    "time_out": "2026-02-26T14:32:10.456Z",
    "duration": 333,
    "return_code": 200,
    "username": "jdoe",
    "userole": "admin",
    "org_id": "org-789",
    "user_agent": "Mozilla/5.0",
    "error_msg": ""
  }'

# Query by username
curl 'http://localhost:8080/httplog/username?username=jdoe'

# Query by URL
curl 'http://localhost:8080/httplog/url?url=/api/v1/orders/12345'
```

