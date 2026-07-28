# Go Clean Architecture Template

A Go template organized as domain modules with DDD and Clean Architecture boundaries. The current example exposes user registration, login, and authenticated profile endpoints.

## Configuration and startup

Configuration is loaded from `.env` when present, then process environment variables take precedence. Copy the sample for local development:

```sh
cp .env.sample .env
# Generate a distinct production secret; do not commit it.
python3 -c 'import secrets; print(secrets.token_urlsafe(48))'
```

`SECRET_KEY` is required and must be at least 32 bytes. Set it through the deployment environment or secret manager; `.env.sample` deliberately contains only a non-secret placeholder. Database values (`DB_HOST`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_PORT`) are also required. `PORT` defaults to `8000`; access and refresh TTLs default to `20m` and `168h`.

Migrations run automatically during application startup. An existing legacy `users` table is accepted as the version-1 baseline only when its required unique email and username constraints are present; otherwise startup fails rather than guessing at schema safety.

The liveness endpoint is `GET /healthz` and returns `200` once the HTTP service is available.

## Run with Docker Compose

Compose requires a real secret in the caller environment; it does not load a secret from the image:

```sh
export SECRET_KEY="$(python3 -c 'import secrets; print(secrets.token_urlsafe(48))')"
docker compose up --build
```

The API is published on `http://localhost:8000` by default. `PORT` controls the port inside the container and `HOST_PORT` controls the published local port; for example, `PORT=8080 HOST_PORT=18080 docker compose up --build` publishes `localhost:18080` to the API listening on `8080`. Compose uses PostgreSQL 16 and persists its data in the `postgres-data` volume.

## API behavior

- `POST /api/auth/signup` creates an account.
- `POST /api/auth/login` returns an access token and a refresh token.
- `GET /api/user/me` requires `Authorization: Bearer <access token>`; refresh tokens are intentionally rejected with `401`.
- Public errors use `{ "code": "...", "message": "..." }`. Generated Swagger documents are in `cmd/api/docs/`.
- Authentication endpoints enforce request-body limits and an in-memory, per-IP rate limit. This limiter is single-process only; use shared edge/distributed controls when running multiple instances.

The HTTP server has bounded read/write timeouts and handles `SIGINT`/`SIGTERM` with graceful shutdown before closing the database connection.

## Development

Prerequisites: the Go version declared in `go.mod`, Docker Compose, Make, and the tools installed by `make init`.

```sh
make doc       # regenerate Swagger artifacts
make test      # test with coverage
make build     # build bin/go-clean-arch
make lint      # run golangci-lint
make e2e       # Compose black-box signup/login/access/refresh smoke test
```

`make e2e` creates a temporary strong `SECRET_KEY`, checks the resolved Compose `HOST_PORT` → `PORT` mapping, starts reusable Compose services with `docker compose up --build`, waits for `/healthz`, validates signup, login, access-token authorization and refresh-token rejection, then removes its Compose resources. Set `PORT` and/or `HOST_PORT` to validate a non-default mapping.

## Extending the application

New HTTP modules implement the shared `RouteRegistrar` contract and are supplied through DI. Add the registrar to the DI provider set and regenerate Wire (`make di`); do not hand-edit central route registration as an extension mechanism. Keep domain, use case, and adapter dependencies pointing inward.

## Project layout

```text
internal/{domain}/domain/   # entities and business rules
internal/{domain}/usecase/  # application workflows and ports
internal/{domain}/adapter/  # HTTP, persistence, and infrastructure adapters
internal/shared/            # shared server, routing, and middleware
internal/di/                # Wire composition root
pkg/                         # configuration, auth, database, logging, utilities
```

## License

MIT. See [LICENSE](LICENSE).
