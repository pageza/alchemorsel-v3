# Alchemorsel v3

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

The server exposes `/api/v1/health` for basic health checks.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Then open `http://localhost:5173`.

## Documentation

Please refer to the markdown files such as `system-overview.md` and `implementation-order.md` for a detailed description of the architecture, API design, and testing approach.

