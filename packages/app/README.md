# App

Vite + React 19 frontend for the pack calculator, styled with Tailwind CSS v4 to match the Gymshark storefront: white surfaces, black type, bold uppercase headings and square black buttons.

The page has two cards: one lists the pack sizes, where you can add a size or remove one with ×, saving each change through `PUT /api/packs`. "Reset to defaults" calls `POST /api/packs/reset`. The other calculates the packs for an order.

## Run

Requires Node 24+ and pnpm 11. Start the server first.

```sh
pnpm install
pnpm dev
```

Open http://localhost:5173. The dev server proxies `/api` to `http://localhost:8080`; in Docker, nginx does the same (`nginx.conf`).

## Scripts

| Script       | Does                          |
| ------------ | ----------------------------- |
| `pnpm dev`   | Dev server with hot reload    |
| `pnpm build` | Type-check and build to dist/ |
| `pnpm test`  | Vitest + Testing Library      |
| `pnpm lint`  | oxlint                        |

## Layout

```
src/
  components/   Reusable UI: Button, Card, DataTable, Input, Spinner, Tag
  layouts/      Page chrome (PlatformLayout)
  pages/        Screens (PackCalculator)
  lib/utils.ts  cn() for merging Tailwind classes
```

Imports use the `@/` alias for `src/`.

## Styling

Tailwind classes only: no CSS files per component and no `style` props. Use `cn()` from `@/lib/utils` when a component accepts a `className`, so callers can override defaults without class conflicts.

## DataTable

Generic table driven by a column config.

```tsx
const columns: DataTableColumn<OrderPack>[] = [
  { key: "packSize", header: "Pack size", cell: (row) => row.packSize },
  { key: "quantity", header: "Quantity", cell: (row) => row.quantity, align: "right" },
];

<DataTable columns={columns} rows={order.packs} rowKey={(row) => row.packSize} caption="Packs to ship" />
```

| Prop           | Type                       | Description                                  |
| -------------- | -------------------------- | -------------------------------------------- |
| `columns`      | `DataTableColumn<TRow>[]`  | `key`, `header`, `cell(row)`, optional `align` |
| `rows`         | `TRow[]`                   | Data to render                               |
| `rowKey`       | `(row) => Key`             | Stable React key per row                     |
| `caption`      | `string`                   | Screen-reader caption                        |
| `emptyMessage` | `ReactNode`                | Shown when `rows` is empty (default "No results") |

## API responses

`request()` in `PackCalculator` trusts the Go API's response shape and does no runtime validation. Add a schema check if the API is ever third-party.
