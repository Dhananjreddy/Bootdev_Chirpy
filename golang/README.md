# 🐦 Chirpy (Golang Backend)

A microblogging API built with Go, designed for learning and experimentation with backend concepts.  
It allows users to register, log in, post short "chirps", and manage authentication securely using JWTs.  
The project uses PostgreSQL for persistent data and `sqlc` for type-safe query generation.

---

## 🚀 Overview

Chirpy is a backend REST API inspired by Twitter’s core idea — short posts ("chirps").  
It’s built with Go’s standard library for high performance and simplicity, emphasizing:
- Clean architecture (`internal` packages)
- Secure authentication (JWT + refresh)
- SQLC for safe database access
- Readiness probes and development environments

---

## 🧰 Tech Stack

| Component | Technology |
|------------|-------------|
| **Language** | Go 1.21+ |
| **Database** | PostgreSQL |
| **ORM/Queries** | SQLC |
| **Auth** | JWT + Refresh Tokens |
| **Config** | Environment Variables |
| **Web** | net/http (standard library) |

---

## 📂 Folder Structure

```
golang/
├── assets/             # Static assets (e.g. logo)
├── chirps.go           # Chirp CRUD endpoints
├── users.go            # User registration/update
├── login.go            # Login & token creation
├── refresh.go          # Refresh token logic
├── reset.go            # Password reset flow
├── readiness.go        # Health & readiness endpoint
├── json.go             # Helper for JSON responses
├── internal/
│   ├── auth/           # JWT, password hashing, secrets
│   └── database/       # DB queries and sqlc integration
├── sql/
│   ├── schema/         # SQL schema
│   └── queries/        # SQLC queries
├── sqlc.yaml           # SQLC configuration
├── main.go             # App entry point
├── go.mod / go.sum     # Dependencies
└── index.html          # Basic test UI
```

---

## ⚙️ Configuration

Example `.env` for development:

```
DB_URL="<your-postgres-access-string-here>"
PLATFORM="dev"
SECRET="<your-jwt-secret-here>"
POLKA_KEY="<your-polka-key-here>"
```

---

## 🧩 Running Locally

1. Clone the repository  
   ```bash
   git clone https://github.com/Dhananjreddy/Bootdev_Chirpy.git
   cd Bootdev_Chirpy/golang
   ```

2. Install dependencies  
   ```bash
   go mod tidy
   ```

3. Start PostgreSQL and ensure `DB_URL` is set correctly.

4. Build and run:  
   ```bash
   go build -o chirpy-server .
   ./chirpy-server
   ```

The API will run at `http://localhost:8080`.

---

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|-----------|-------------|
| `POST` | `/api/users` | Register a new user |
| `POST` | `/api/login` | Authenticate user |
| `POST` | `/api/refresh` | Refresh access token |
| `POST` | `/api/reset` | Reset password |
| `GET`  | `/api/chirps` | List all chirps |
| `GET`  | `/api/chirps/{id}` | Get single chirp |
| `POST` | `/api/chirps` | Create chirp (auth) |
| `DELETE` | `/api/chirps/{id}` | Delete chirp (owner only) |
| `GET`  | `/api/readiness` | Health check |

---

## 🧪 Testing

Run all tests:
```bash
go test ./...
```

---
