import { render, screen, within } from "@testing-library/react";
import { DataTable, type DataTableColumn } from "./DataTable";

interface Row {
  id: number;
  name: string;
}

const columns: DataTableColumn<Row>[] = [
  { key: "id", header: "ID", cell: (row) => row.id, align: "right" },
  { key: "name", header: "Name", cell: (row) => row.name },
];

describe("DataTable", () => {
  it("renders a header and a row per item", () => {
    render(<DataTable columns={columns} rows={[{ id: 1, name: "Shorts" }, { id: 2, name: "Hoodie" }]} rowKey={(row) => row.id} />);

    const [header, ...rows] = screen.getAllByRole("row");
    expect(within(header).getAllByRole("columnheader").map((cell) => cell.textContent)).toEqual(["ID", "Name"]);
    expect(rows).toHaveLength(2);
    expect(within(rows[1]).getByText("Hoodie")).toBeInTheDocument();
    expect(within(rows[0]).getByText("1")).toHaveClass("text-right");
  });

  it("shows the empty message when there are no rows", () => {
    render(<DataTable columns={columns} rows={[]} rowKey={(row) => row.id} emptyMessage="Nothing here" />);

    expect(screen.getByRole("cell")).toHaveTextContent("Nothing here");
  });
});
