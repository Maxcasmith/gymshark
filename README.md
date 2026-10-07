# Gymshark Pack Calculator

Works out which whole packs to ship for an order. The rules, in priority order:

1. Only whole packs are sent.
2. Ship the fewest items that still cover the order.
3. Within rule 2, ship as few packs as possible.

Pack sizes can be added and removed in the UI (or via `PUT /api/packs`); they live in `db/packs.json`, so no code changes are needed.

```
db/                 JSON data (packs.json, orders.json)
packages/server     Go API
packages/app        Vite + React frontend
```

## Run with Docker

Requires Docker with Compose.

```sh
docker compose up --build
```

- App: http://localhost:5173
- API: http://localhost:8080

Orders are stored in the `data` volume and survive restarts. The volume is seeded from `db/` the first time it is created, so after editing `db/packs.json` reset it:

```sh
docker compose down -v
docker compose up --build
```

## Run locally

See each package's README:

- [packages/server](packages/server/README.md): Go API, flags, API reference, algorithm notes
- [packages/app](packages/app/README.md): React frontend, components, styling

## Tests

```sh
cd packages/server && go test ./...
cd packages/app && pnpm test
```
