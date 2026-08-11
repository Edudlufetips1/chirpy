# Chirpy

Chirpy is a small Go API for a simple microblogging app. It handles user accounts, JWT auth, chirp creation, chirp retrieval, sorting, deletion, and a basic webhook flow for upgrading users.

This project uses PostgreSQL for storage and `sqlc` to generate type-safe database access code.

## What it does

- Create users
- Log in and authenticate requests with JWTs
- Refresh and revoke access tokens
- Create chirps with profanity filtering
- List chirps, sort them, and filter by author
- Fetch a single chirp by ID
- Delete only your own chirps
- Handle Polka webhook upgrades for premium users
- Reset the app state in development mode

## Prerequisites

- Go 1.22+
- PostgreSQL running locally
- `sqlc` installed
- A database created for the app

## Setup

1. Create a `.env` file in the project root:

   ```env
   DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
   JWT_SECRET=your_jwt_secret
   PLATFORM=dev
   POLKA_KEY=your_polka_api_key
   ```

2. Make sure your PostgreSQL database is running and accessible with the `DB_URL` above.

3. Generate the typed database code:

   ```bash
   sqlc generate
   ```

4. Start the server:

   ```bash
   go run .
   ```

5. The API starts on port `8080`.

## Project layout

- `main.go` — app setup and routing
- `handler_*.go` — HTTP handlers for users, auth, chirps, and webhooks
- `internal/auth/` — JWT and API key helpers
- `internal/database/` — generated database queries
- `sql/schema/` — database schema files
- `sql/queries/` — SQL queries used by `sqlc`

## API overview

### Users

- `POST /api/users`
  - Create a new user
  - Body: `{ "email": "user@example.com", "password": "secret" }`

### Auth

- `POST /api/login`
  - Logs in a user and returns JWT + refresh token

- `POST /api/refresh`
  - Refreshes an access token using a refresh token

- `POST /api/revoke`
  - Revokes a refresh token

### Chirps

- `POST /api/chirps`
  - Creates a chirp
  - Requires a valid bearer JWT
  - Body: `{ "body": "hello world" }`

- `GET /api/chirps`
  - Returns all chirps
  - Optional query param: `?sort=asc` or `?sort=desc`

- `GET /api/chirps/{chirpID}`
  - Returns one chirp by ID

- `DELETE /api/chirps/{chirpID}`
  - Deletes a chirp
  - Requires a valid bearer JWT
  - Only the chirp owner can delete it

### Webhooks

- `POST /api/polka/webhooks`
  - Validates the `Authorization: ApiKey ...` header
  - Confirms the event is `user.upgraded`
  - Upgrades the matched user to Chirpy Red

### Admin

- `GET /admin/metrics`
  - Shows app hit count

- `POST /admin/reset`
  - Clears hit count in dev mode only

## Auth notes

The app expects bearer tokens for normal user actions:

```http
Authorization: Bearer <jwt>
```

Webhook requests use an API key header:

```http
Authorization: ApiKey <polka_key>
```

## Notes

- Chirps are checked for profanity before they are saved.
- `sqlc` generates the database code, so you generally do not hand-write query functions.
- The app is designed to be easy to run locally for development.
