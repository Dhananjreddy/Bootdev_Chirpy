# Chirpy Backend

A lightweight backend API for Chirpy, a microblogging platform. Handles users, chirps, authentication, and membership status. Hit the project with a star if you find it useful ⭐

Supported by Boot.dev.

---

## 🚀 Deploy

Clone the repository, install dependencies, and run the server:

```bash
git clone <repo-url>
cd typescript
npm install
npm run migrate
npm run dev
```

Server runs at `http://localhost:8080`.

---

## 💡 Motivation

Chirpy's backend is designed to provide a simple, maintainable, and fully typed API for a Twitter-like experience. It leverages TypeScript and Drizzle ORM for type safety, PostgreSQL for storage, and JWT/refresh token authentication for security.

Goals:

* Full CRUD for chirps and users
* Authentication and refresh token handling
* Membership management (isChirpyRed)
* Readiness and metrics endpoints for monitoring

---

## ⚙️ Installation

Inside the project folder:

```bash
npm install
```

Run migrations and start the server:

```bash
npm run migrate
npm run dev
```

---

## 📝 Quick Start

### Users

* **Create User**: `POST /api/users`
  Body: `{ "email": "user@example.com", "password": "password123" }`

* **Update User**: `PUT /api/users`
  Body: `{ "email": "newemail@example.com" }`

### Chirps

* **Create Chirp**: `POST /api/chirps`
  Body: `{ "body": "Hello Chirpy!" }`

* **Get All Chirps**: `GET /api/chirps`

* **Get Chirp by ID**: `GET /api/chirps/:chirpID`

* **Delete Chirp**: `DELETE /api/chirps/:chirpID`

### Authentication

* **Login**: `POST /api/login`
  Body: `{ "email": "user@example.com", "password": "password123" }`

* **Refresh Token**: `POST /api/refresh`
  Body: `{ "refreshToken": "<token>" }`

* **Revoke Token**: `POST /api/revoke`
  Body: `{ "refreshToken": "<token>" }`

### Membership

* **Upgrade Membership**: `POST /api/polka/webhooks`
  Webhook handler to set `isChirpyRed` on user upgrade events.

### Admin

* **Reset**: `POST /admin/reset`
* **Metrics**: `GET /admin/metrics`

### Health

* **Readiness**: `GET /api/healthz`

---

## 🗂️ Project Structure

```
src/
├─ api/
│  ├─ createUser.ts
│  ├─ createChirp.ts
│  ├─ getChirps.ts
│  ├─ login.ts
│  ├─ membership.ts
│  ├─ validate.ts
│  ├─ readiness.ts
│  └─ middleware.ts
├─ db/
│  ├─ schema.ts
│  ├─ queries/
│  └─ migrations/
├─ auth.ts
├─ index.ts
└─ server.ts
```

---

## 📦 Dependencies

* Express
* Drizzle ORM
* PostgreSQL
* Argon2 for password hashing
* JSON Web Tokens (JWT)
