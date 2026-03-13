# Email Worker Service

A microservice for rendering and delivering transactional emails, triggered by the Risefit API.

## Features
- **Framework**: Gin (Web) + Google Wire (Dependency Injection).
- **Endpoint**: `POST /send-email`
- **Authentication**: `X-Internal-API-Key` header validation.
- **Templates**: `verify_email`, `reset_password`.
- **ESP**: SendGrid (configurable).
- **Retries**: Returns `5xx` for transient failures, enabling Cloud Tasks retries.

## API Documentation

### Send Email
**Endpoint**: `POST /send-email`  
**Headers**: `X-Internal-API-Key: <your_internal_key>`

#### Request Payload
```json
{
  "email": "recipient@example.com",
  "template": "template_name",
  "data": {
    "key": "value"
  }
}
```

#### Supported Templates
| Template Name | Required `data` Fields | Description |
| :--- | :--- | :--- |
| `verify_email` | `name`, `verification_url` | New user account confirmation. |
| `reset_password` | `name`, `reset_url` | Password recovery link. |

#### Example Payloads
**Verify Email**:
```json
{
  "email": "user@example.com",
  "template": "verify_email",
  "data": {
    "name": "John Doe",
    "verification_url": "https://risefit.com/verify?token=abc"
  }
}
```

**Reset Password**:
```json
{
  "email": "user@example.com",
  "template": "reset_password",
  "data": {
    "name": "John Doe",
    "reset_url": "https://risefit.com/reset?token=123"
  }
}
```

## Environment Variables
- `PORT`: Port to listen on (default: `8080`).
- `INTERNAL_API_KEY`: Secret key shared with the main API.
- `EMAIL_PROVIDER_API_KEY`: SendGrid API key.
- `FROM_EMAIL`: Sender address (e.g., `no-reply@risefit.com`).

## Setup & Running

### 1. Get a SendGrid API Key
1.  Log in to [SendGrid](https://app.sendgrid.com/).
2.  Go to **Settings > API Keys** and create a key with **Mail Send** permissions.
3.  Verify your sender identity in **Settings > Sender Authentication**.

### 2. Running Locally
```bash
export PORT=8080
export INTERNAL_API_KEY=your_internal_key
export EMAIL_PROVIDER_API_KEY=your_sendgrid_key
export FROM_EMAIL=no-reply@risefit.com

go run cmd/worker/*.go
```

### 3. Running with Docker
```bash
docker build -t email-worker .
docker run -p 8080:8080 \
  -e INTERNAL_API_KEY=your_internal_key \
  -e EMAIL_PROVIDER_API_KEY=your_sendgrid_key \
  -e FROM_EMAIL=no-reply@risefit.com \
  email-worker
```

## Development
This project uses **Google Wire** for dependency injection. If you modify dependencies in `cmd/worker/wire.go`, regenerate the injector:
```bash
cd cmd/worker
go run github.com/google/wire/cmd/wire
```
