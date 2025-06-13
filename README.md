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

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Then open `http://localhost:5173`.

## Documentation

Please refer to the markdown files such as `system-overview.md` and `implementation-order.md` for a detailed description of the architecture, API design, and testing approach.

