# AY.com – Fullstack Twitter Clone

AY.com is a **thread-based social media platform** inspired by Twitter/X. It allows users to share short posts (threads), images, and videos, and to engage with each other through likes, replies, reposts, communities, and real-time messaging.

This project was developed as a **fullstack application** with modern web technologies, focusing on scalability, microservices, and security.

---

## 🚀 Tech Stack

### Frontend

* **Framework:** [Svelte](https://svelte.dev/) + [Vite](https://vitejs.dev/)
* **Language:** TypeScript
* **Styling:** Pure SCSS (no UI libraries)
* **Features:**

  * Responsive design (desktop, tablet, mobile)
  * Light/Dark mode toggle
  * Infinite scrolling feeds
  * Rich text rendering (mentions + hashtags)
  * Loading skeletons

### Backend

* **Language:** Go (Golang)
* **Architecture:** Microservices
* **Communication:** gRPC + RabbitMQ
* **API Gateway:** Bridges HTTP frontend ↔ gRPC services
* **API Documentation:** Swagger

### Databases & Caching

* **PostgreSQL** (separate DB per service)
* **Redis** (caching & performance)

### Authentication & Security

* JWT (access + refresh tokens)
* Password salting & hashing
* Google OAuth login
* reCAPTCHA

### Media & AI

* **Supabase** – store images, GIFs, videos
* **Flask** – deploy AI features (speech recognition, search enhancements)

### DevOps & Tooling

* Docker + Docker Compose
* ESLint (provided config)
* Unit testing (Go + mocks)
* Structured logging (configurable log levels)

---

## 📂 Project Structure

```
.
├── backend
│   ├── api-gateway           # HTTP → gRPC gateway
│   ├── proto                 # Generated protobuf files
│   ├── service-ai            # AI-related Flask integrations
│   ├── service-community     # Communities service
│   ├── service-media         # Media storage service
│   ├── service-message       # Messaging service (real-time chat)
│   ├── service-notification  # Notifications service
│   ├── service-thread        # Threads & posts service
│   ├── service-user          # Authentication & user profiles
│   ├── util-email            # Email utility (verification, notifications)
│   └── util-redis            # Redis utility
│
├── frontend
│   ├── public                # Static assets
│   ├── src                   # Svelte + Vite + TS frontend code
│   ├── package.json          # Frontend dependencies
│   └── Dockerfile            # Container config
│
├── compose.yaml              # Alternative Docker config
├── docker-compose.yml        # Main Docker Compose config
├── start-all.bat             # Windows startup script for development
├── .dockerignore
├── .gitignore
└── README.md
```

---

## ⚡ Installation & Setup

### Prerequisites

* Node.js (>=18)
* Go (>=1.21)
* PostgreSQL
* Redis
* Docker & Docker Compose

### Steps

1. **Clone the repository**

   ```bash
   git clone https://github.com/nathabuddhi/ay-com.git
   cd ay-com
   ```

2. **Setup environment variables**
   Create `.env` files for frontend and each backend service.
   Include values for:

   * Database URIs
   * Redis connection
   * JWT secrets
   * Supabase keys
   * RabbitMQ config

3. **Run on Windows**
   Simply execute:

   ```bash
   start-all.bat
   ```

   This script will start all microservices, the API gateway, and the frontend.
   Use docker desktop if you prefer to do so.

5. **Access the app**

   * Frontend: `http://localhost:5173`
   * API Gateway: `http://localhost:8080`
   * Swagger Docs: `http://localhost:8080/docs`

---

## 🔮 Future Improvements

* Push notifications (browser & mobile)
* Advanced AI moderation (spam/abuse detection)
* GraphQL gateway alternative

---
