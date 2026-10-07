# Server

Go API that stores pack sizes and orders in JSON files and works out which packs to ship. Uses the standard library only.

## Run

Requires Go 1.27+.

```sh
go run .
```

| Flag       | Default                 | Description                                  |
| ---------- | ----------------------- | -------------------------------------------- |
| `-addr`    | `:8080`                 | Listen address                               |
| `-db`      | `../../db`              | Directory holding `packs.json`/`orders.json` |
| `-origins` | `http://localhost:5173` | Comma-separated origins allowed by CORS      |

Running locally writes orders to `../../db/orders.json`.

## Test

```sh
go test ./...
```

## Layout

```
server.go                 Flags, CORS, starts the HTTP server
routes.go                 HTTP handlers and error-to-status mapping
entities/                 Pack and Order types
commands/                 Write requests (CreateOrderCommand, ReplacePacksCommand) and their handlers
queries/                  Read requests (GetPackQuery, GetOrderQuery) and their handlers
repositories/             JSON file storage, one repository per entity
helpers/OrderAggregator   Pack calculation
middleware/Cors           Allows the listed origins and answers preflight requests
```

Commands and queries follow a small CQRS split. `commands/Base` and `queries/Base` define the handler interfaces; each handler asserts it satisfies them at compile time.

## Pack calculation

`OrderAggregator` picks whole packs that, in priority order:

1. ship the fewest items that still cover the order
2. use the fewest packs for that item count

Greedy largest-first breaks rule 1 when sizes are not multiples of each other (for example 23/31/53), so it runs an unbounded min-coin dynamic programme over every total.

`ceil(requested / smallest) * smallest` is always reachable and is below `requested + smallest`, so the search only covers `[requested, requested + smallest)`. Two tables track, per total, the fewest packs needed (`0` means unreachable) and the last pack added, which is walked back to build the breakdown.

Cost is O(items × sizes) time and memory, which is why orders are capped at `MaxRequestedItems` (1,000,000). To go higher, divide totals by the GCD of the pack sizes first.

## Storage

Each repository is one JSON file guarded by a mutex, so it is safe within a single process only. `Append` writes to a `.tmp` file and renames it over the original, so a crash mid-write never truncates the data. Swap for a real database if more than one instance must run.

## API

| Method | Path               | Body                        | Response                   |
| ------ | ------------------ | --------------------------- | -------------------------- |
| GET    | `/api/packs`       |                             | `[{ "id", "size" }]`       |
| PUT    | `/api/packs`       | `{ "sizes": [250, 500] }`   | `[{ "id", "size" }]`       |
| POST   | `/api/packs/reset` |                             | `[{ "id", "size" }]`       |
| POST   | `/api/orders`      | `{ "requestedItems": 501 }` | Order (201)                |
| GET    | `/api/orders/{id}` |                             | Order, or 404 if not found |

Order:

```json
{
  "id": 1,
  "requestedItems": 501,
  "shippedItems": 750,
  "excessItems": 249,
  "totalPacks": 2,
  "packs": [
    { "packSize": 500, "quantity": 1 },
    { "packSize": 250, "quantity": 1 }
  ]
}
```

`PUT /api/packs` replaces every pack size. Sizes are sorted, duplicates dropped and ids reassigned from 1. It accepts 1 to 20 sizes, each between 1 and 1,000,000, which keeps the calculation within the same bounds as orders. `POST /api/packs/reset` restores the sizes from the brief (`DefaultPackSizes`: 250, 500, 1000, 2000, 5000). Changes apply to the next order, and editing `packs.json` by hand works too.

Errors return `{ "error": "message" }` with status 400, 404 or 500. Request bodies are limited to 1 KB.

## Docker

The image is built from the repo root so the seed data in `db/` can be copied to `/data`. A named volume mounted at `/data` is initialised from it on first run.
