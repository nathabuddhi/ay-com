# TPA WEB NB24-2

# FIRST COMMIT

## Running with Docker

This project uses Docker Compose to orchestrate all backend, frontend, and supporting services. Each service has its own Dockerfile and is configured in `docker-compose.yml`.

### Requirements
- Docker and Docker Compose installed
- Docker Compose will build images using:
  - Go version: **1.24.2** (for backend services)
  - Node.js version: **22.13.1** (for frontend)

### Environment Variables
Each service loads its environment variables from its own `.env` file. Make sure to fill in the required values in these files before running:
- `./backend/api-gateway/.env`
- `./backend/service-media/.env`
- `./backend/service-notification/.env`
- `./backend/service-thread/.env`
- `./backend/service-user/.env`
- `./backend/util-email/.env`
- `./backend/util-redis/.env`
- `./frontend/.env`

### Build and Run
To build and start all services:

```sh
# From the project root directory
docker compose up --build
```

This will build and start all services as defined in `docker-compose.yml`.

### Exposed Ports
- **API Gateway:** `5000` (HTTP)
- **Service User:** `5001` (gRPC)
- **Service Thread:** `5002` (gRPC)
- **Service Notification:** `5009` (gRPC)
- **Service Media:** `5013` (gRPC)
- **Frontend (Svelte/Vite):** `5173` (HTTP)
- **RabbitMQ:** `5672` (AMQP), `15672` (Management UI)
- **Redis:** `6379`

### Special Configuration
- All services are connected via the `backend` Docker network.
- RabbitMQ and Redis are included as containers and can be accessed by the backend services.
- To persist RabbitMQ or Redis data, uncomment the `volumes` sections in `docker-compose.yml`.
- All Go services are built as static binaries and run as non-root users for security.
- The frontend is served using Vite's preview server.

---
