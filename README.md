# Alchemorsel v3

 contains documentation and a project skeleton for **Alchemorsel**, an AI‑powered recipe application.  The code aligns with the architecture and API design documents shipped in this repository.
=======
This repository contains documentation and a small project skeleton for **Alchemorsel**, an AI‑powered recipe application. The architecture is based on the markdown files found in this repo.


## Contents

- `backend` &ndash; Go API server using Gin
- `frontend` &ndash; Vue 3 app built with Vite
- Design documents in markdown and diagrams

## Usage

### Backend

```bash
cd backend
go run ./cmd/api
```


The server exposes a versioned API under `/api/v1`. Only a simple health check is implemented but the router already includes endpoints for authentication, user profiles and recipes as described in the design docs.
The router also applies placeholder middleware for authentication, CORS, logging, and panic recovery.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Then open `http://localhost:5173`.

### Configuration

Copy `.env.example` to `.env` and adjust any of the values if needed. The
backend reads these variables when starting up.

### Makefile

Common commands are wrapped in a simple `Makefile`. From the project root you
can run:

```bash
make run  # start the API server
make test # run backend tests
```

### Git Hooks

Run the following command once to enable the provided hooks:

```bash
git config core.hooksPath .githooks
```

The `pre-commit` hook verifies code formatting and runs tests before allowing a
commit.

### Docker Compose

`docker-compose.yml` launches the API together with PostgreSQL and Redis. Start
all services with:

```bash
docker-compose up
```

## Documentation

Please refer to the markdown files such as `system-overview.md` and `implementation-order.md` for a detailed description of the architecture, API design, and testing approach.

