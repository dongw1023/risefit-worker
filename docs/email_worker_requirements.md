# Requirements: Email Worker Service

This document specifies the requirements for the **Email Worker Service**, a dedicated microservice responsible for rendering and delivering transactional emails triggered by the Risefit API.

---

## 1. Functional Requirements

### 1.1 Task Processing
- **Endpoint**: `POST /send-email`
- **Input**: A JSON payload from Google Cloud Tasks (see `docs/worker_infrastructure.md`).
- **Templates**: Support HTML email templates with dynamic data injection.
- **Delivery**: Integrate with a Transactional Email Provider (ESP) like SendGrid, Postmark, or AWS SES.

### 1.2 Supported Email Templates
The service must support the following templates at launch:
1.  **`verify_email`**:
    - **Purpose**: New user account confirmation.
    - **Required Data**: `name`, `verification_url`.
2.  **`reset_password`**:
    - **Purpose**: Password recovery link.
    - **Required Data**: `name`, `reset_url`.

---

## 2. Technical Specifications

### 2.1 Technology Stack (Recommended)
- **Language**: Go (to match the main API) or Node.js (excellent for template rendering).
- **Template Engine**: 
    - Go: `html/template`
    - Node.js: `handlebars` or `ejs`
- **Email SDK**: Use the official SDK of the chosen ESP (e.g., `@sendgrid/mail`).

### 2.2 Security
- **Authentication**: Every request must be validated using the `X-Internal-API-Key` header. Requests without a valid key must return `401 Unauthorized`.
- **Validation**: Validate that the `email` field is a valid email address and all required fields for the specific `template` are present in the `data` object.

### 2.3 Error Handling & Retries
- **Retry Logic**: If the ESP returns a temporary error (e.g., 5xx, Rate Limit), the worker must return a `5xx` status code. This signals Google Cloud Tasks to retry the task automatically based on the queue's backoff policy.
- **Fatal Errors**: If the request is invalid (e.g., 4xx), return a `4xx` status code to Cloud Tasks so it does **not** retry a permanently broken task.

---

## 3. Operational Requirements

### 3.1 Logging
- Log every incoming task with its `type` and recipient `email` (do not log sensitive tokens).
- Log any delivery failures with the error message from the ESP.

### 3.2 Deployment
- **Platform**: Google Cloud Run (recommended) or any containerized environment.
- **Environment Variables**:
    - `PORT`: Port the service listens on.
    - `INTERNAL_API_KEY`: Secret key to match the main API.
    - `EMAIL_PROVIDER_API_KEY`: API key for SendGrid/Postmark/etc.
    - `FROM_EMAIL`: The "From" address (e.g., `no-reply@risefit.com`).

---

## 4. Example Workflow
1.  **API**: User signs up. API sends a task to `GCP_EMAIL_QUEUE_NAME`.
2.  **Cloud Tasks**: Delivers the task to the Worker's `POST /send-email` endpoint.
3.  **Worker**: 
    - Verifies `X-Internal-API-Key`.
    - Matches `template: "verify_email"`.
    - Renders `verify_email.html` using `data.name` and `data.verification_url`.
    - Calls SendGrid API to send the final HTML.
    - Returns `200 OK` to Cloud Tasks.
