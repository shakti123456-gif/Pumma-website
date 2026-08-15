# Ecommerce store

## Storefront

```bash
cd frontend
npm install
npm run dev
```

## Backend (Gin)

```bash
cd backend
go mod tidy
cp .env.example .env
# Set DATABASE_URL in .env, then export it (or use your deployment's environment variables)
go run .
```

The API runs at `http://localhost:8080`. It requires `DATABASE_URL` and creates the CockroachDB-compatible tables on startup.

## Admin panel

In a second terminal:

```bash
cd admin
npm install
npm run dev
```

Open the URL printed by Vite (normally `http://localhost:5174`). Set `VITE_API_URL` in an `admin/.env` file when the API is hosted elsewhere.

## Important

This is a development foundation. Add authentication, a database, image uploads, server-side validation, and environment-restricted CORS before a production launch.
