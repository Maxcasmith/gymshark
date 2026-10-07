import type { Key, ReactNode } from "react";
import { cn } from "@/lib/utils";

export interface DataTableColumn<TRow> {
  key: string;
  header: ReactNode;
  cell: (row: TRow) => ReactNode;
  align?: "left" | "right";
}

interface DataTableProps<TRow> {
  columns: DataTableColumn<TRow>[];
  rows: TRow[];
  rowKey: (row: TRow) => Key;
  caption?: string;
  emptyMessage?: ReactNode;
}

export function DataTable<TRow>(props: DataTableProps<TRow>) {
  const { columns, rows, rowKey, caption, emptyMessage = "No results" } = props;

  return (
    <table className="w-full text-left text-sm">
      {caption && <caption className="sr-only">{caption}</caption>}
      <thead>
        <tr className="border-b border-black">
          {columns.map((column) => (
            <th
              key={column.key}
              scope="col"
              className={cn(
                "py-3 text-xs font-bold uppercase tracking-wider text-black",
                column.align === "right" && "text-right",
              )}
            >
              {column.header}
            </th>
          ))}
        </tr>
      </thead>
      <tbody>
        {rows.length ? (
          rows.map((row) => (
            <tr key={rowKey(row)} className="border-b border-neutral-200 last:border-b-0">
              {columns.map((column) => (
                <td key={column.key} className={cn("py-3", column.align === "right" && "text-right")}>
                  {column.cell(row)}
                </td>
              ))}
            </tr>
          ))
        ) : (
          <tr>
            <td colSpan={columns.length} className="py-6 text-center text-neutral-500">
              {emptyMessage}
            </td>
          </tr>
        )}
      </tbody>
    </table>
  );
}
