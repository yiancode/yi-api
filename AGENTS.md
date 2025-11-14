# Repository Guidelines

## Project Structure & Module Organization
The Go backend sits at repo root: `main.go` wires `router/`, middleware in `middleware/`, controllers in `controller/`, domain logic in `service/`, and DTO/model types under `dto/`, `model/`, and `types/`. Shared helpers live in `common/`, `constant/`, `logger/`, and `setting/`, while the React/Vite dashboard resides in `web/` (`web/public/` assets, `web/src/components` UI, `web/src/i18n/` locale data). Deployment aides stay in `bin/`, `docker-compose.yml`, `new-api.service`, and `docs/`.

## Build, Test & Development Commands
- `make build-frontend`: Bun install + Vite build with the `VERSION` stamp.
- `make start-backend`: `go run main.go`; stop via `Ctrl+C` once the background process spawns.
- `bun run dev --cwd web`: Hot-reloads the dashboard proxied to `http://localhost:3000`.
- `bun run build --cwd web`: Production build smoke test; required before UI pull requests.
- `go test ./...`: Runs all Go tests; narrow scope with `./service/... -run Foo`.
- `docker-compose up -d`: Spins the full stack (SQLite default); override env secrets in compose.

## Coding Style & Naming Conventions
Format Go code with `gofmt`/`goimports` (tabs, PascalCase exports, snake_case only when mirroring SQL columns). Keep controllers thin, move integrations into `service/`, and wrap errors with `fmt.Errorf("context: %w", err)`. Frontend files honor the repo Prettier single-quote config plus ESLint React rules; run `bun run lint` and `bunx eslint "**/*.{js,jsx}" --cache`. Components belong in `web/src/components` or `web/src/pages`, and custom hooks must start with `use`.

## Testing Guidelines
Store Go specs beside their code as `*_test.go` table-driven cases; mock providers with in-memory stubs and aim for >80% coverage on new services using `go test -cover ./service/... ./controller/...`. The web layer currently relies on lint/build feedback; optional Vitest/Jest specs live in `web/src/__tests__` and should stub `/api`.

## Commit & Pull Request Guidelines
Recent commits use short imperative subjects (`fix GetChannelKey...`, `support replicate channel`); keep them under 72 characters, mention the module when relevant, and squash noisy WIP commits. PRs should spell out motivation, link issues, flag config or migration changes, attach UI screenshots, and confirm `go test ./...` plus `bun run build --cwd web`.

## Security & Configuration Tips
`setting/config/config.go` sources runtime options, so inject secrets via environment variables or compose overrides rather than code. Persist SQLite data with the `./data` volume defined in `docker-compose.yml`. Rotate tokens referenced in `setting/*.go` if logs expose them, and ensure `.env*` files stay untracked.
