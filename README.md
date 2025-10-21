# Chirpy

Chirpy is a full-featured microblogging API inspired by Twitter, implemented in both **Golang** and **TypeScript**. It supports user registration, authentication, and CRUD operations for short posts (“chirps”), with secure session management and scalable architecture.  

The **Golang backend** uses Go’s standard `net/http` library, modular internal packages, and **PostgreSQL** integration via **SQLC** for type-safe database queries. Authentication is implemented with **JWT** and rotating refresh tokens, combined with **Argon2id** password hashing for enhanced security.  

The **TypeScript backend** is built with **Express** and **Drizzle ORM**, leveraging strong typing, modular routing, and seamless database migrations. Both backends provide readiness probes, health checks, and structured error handling to support production-ready deployment.  

Chirpy emphasizes **clean architecture, maintainability, and cross-stack consistency**, with environment-based configuration, testing pipelines, and easy setup for local development. It serves as a learning-focused platform for developers exploring secure backend design, RESTful API development, and type-safe database access in Go and TypeScript.  

Quickly start Chirpy by cloning the repo, installing dependencies, running migrations, and starting the server.