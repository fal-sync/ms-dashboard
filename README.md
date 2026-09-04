# Go Clean Architecture Boilerplate

A skeleton for building API services that stay organized as the codebase grows with additional databases, microservice integrations, or event brokers like Kafka.

## Project Structure

```text
.
├── cmd/api                                      # application entry point
├── internal/app                                 # dependency wiring and integration bootstrap
├── internal/config                              # environment config loader
├── internal/delivery/http                       # HTTP handlers and router
├── internal/domain                              # entities, domain errors, domain events
├── internal/usecase                             # business logic
├── internal/usecase/port                        # outbound ports for external integrations
├── internal/infrastructure/persistence          # database / repository adapters
├── internal/infrastructure/external/httpclient  # shared HTTP client for service-to-service calls
├── internal/infrastructure/external/service     # adapters for calling external services
└── internal/infrastructure/external/messaging   # event publishers (e.g. Kafka)
```

## Endpoints

| Method | Path          | Description            |
|--------|---------------|------------------------|
| GET    | /health       | Health check           |
| POST   | /users        | Create a new user      |
| GET    | /users        | List all users         |
| GET    | /users/{id}   | Get a user by ID       |

Example request:

```bash
curl --location 'http://localhost:8080/users' \
  --header 'Content-Type: application/json' \
  --data '{
    "name": "Naufal",
    "email": "naufal@example.com"
  }'
```

## Getting Started

1. Copy the example environment file and adjust values as needed:

   ```bash
   cp .env.example .env
   ```

2. Run the application:

   ```bash
   go run ./cmd/api
   ```

3. Run the tests:

   ```bash
   go test ./...
   ```

## External Integrations

The following integrations are pre-wired and opt-in via environment variables:

| Integration      | Trigger condition                                   | Adapter file |
|------------------|-----------------------------------------------------|--------------|
| HTTP sync hook   | `EXTERNAL_USER_SYNC_BASE_URL` is set (non-empty)   | [`internal/infrastructure/external/service/user_sync_hook.go`](internal/infrastructure/external/service/user_sync_hook.go) |
| Kafka publisher  | `KAFKA_ENABLED=true` **and** brokers are configured | [`internal/infrastructure/external/messaging/kafka/user_created_hook.go`](internal/infrastructure/external/messaging/kafka/user_created_hook.go) |

The use case depends only on the port interface defined in [`internal/usecase/port/user_created_hook.go`](internal/usecase/port/user_created_hook.go), not on any concrete implementation.

## Environment Variables

Copy [`.env.example`](.env.example) to `.env` and fill in your values. Full reference:

| Variable                       | Default                          | Description                                         |
|--------------------------------|----------------------------------|-----------------------------------------------------|
| `APP_NAME`                     | `clean-architecture-boilerplate` | Service name                                        |
| `APP_PORT`                     | `8080`                           | HTTP listener port                                  |
| `APP_SHUTDOWN_TIMEOUT`         | `10s`                            | Graceful shutdown timeout                           |
| `EXTERNAL_USER_SYNC_BASE_URL`  | *(empty — disabled)*             | Base URL of the external user-sync service          |
| `EXTERNAL_USER_SYNC_PATH`      | `/internal/users/sync`           | Request path for the sync endpoint                  |
| `EXTERNAL_USER_SYNC_TIMEOUT`   | `3s`                             | Timeout for the sync HTTP call                      |
| `KAFKA_ENABLED`                | `false`                          | Set to `true` to enable Kafka event publishing      |
| `KAFKA_BROKERS`                | `localhost:9092`                 | Comma-separated list of Kafka broker addresses      |
| `KAFKA_USER_CREATED_TOPIC`     | `user.created`                   | Topic to publish `user.created` events to           |

## Architecture Notes

- **Outbox pattern**: For truly mission-critical event publishing, Kafka writes should go through an outbox pattern rather than being published directly inside the request cycle.
- **Hook placement**: External hooks are called after the user has been persisted, keeping the dependency structure clear and easy to extend.
- **In-memory repository**: The persistence layer currently uses an in-memory store. Swapping it for a real database only requires replacing the repository adapter — the use case and delivery layers are unaffected.
