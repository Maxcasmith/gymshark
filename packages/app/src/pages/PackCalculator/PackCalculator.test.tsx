import userEvent from "@testing-library/user-event";
import { render, screen, waitFor } from "@testing-library/react";
import { PackCalculator } from "./PackCalculator";

function jsonResponse(body: object, status = 200) {
  return new Response(JSON.stringify(body), { status });
}

describe("PackCalculator", () => {
  afterEach(() => vi.unstubAllGlobals());

  it("shows the pack breakdown for a submitted order", async () => {
    const user = userEvent.setup();
    vi.stubGlobal(
      "fetch",
      vi.fn((path: string) =>
        Promise.resolve(
          path === "/api/packs"
            ? jsonResponse([{ id: 1, size: 250 }, { id: 2, size: 500 }])
            : jsonResponse({
                id: 1,
                requestedItems: 501,
                shippedItems: 750,
                excessItems: 249,
                totalPacks: 2,
                packs: [
                  { packSize: 500, quantity: 1 },
                  { packSize: 250, quantity: 1 },
                ],
              }),
        ),
      ),
    );

    render(<PackCalculator />);
    await user.type(await screen.findByLabelText("Items ordered"), "501");
    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(await screen.findByText("Total packs")).toBeInTheDocument();
    expect(screen.getByText("750")).toBeInTheDocument();
    expect(screen.getByText("249")).toBeInTheDocument();
  });

  it("shows the API error message", async () => {
    const user = userEvent.setup();
    vi.stubGlobal(
      "fetch",
      vi.fn((path: string) =>
        Promise.resolve(
          path === "/api/packs"
            ? jsonResponse([{ id: 1, size: 250 }])
            : jsonResponse({ error: "invalid order: items must be between 1 and 1000000" }, 400),
        ),
      ),
    );

    render(<PackCalculator />);
    await user.type(await screen.findByLabelText("Items ordered"), "2000000");
    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("items must be between 1 and 1000000");
  });

  it("adds and removes pack sizes", async () => {
    const user = userEvent.setup();
    const fetchMock = vi.fn((_path: string, init?: RequestInit) => {
      if (init?.method !== "PUT") {
        return Promise.resolve(jsonResponse([{ id: 1, size: 250 }, { id: 2, size: 500 }]));
      }
      const { sizes } = JSON.parse(String(init.body)) as { sizes: number[] };
      return Promise.resolve(jsonResponse(sizes.map((size, index) => ({ id: index + 1, size }))));
    });
    vi.stubGlobal("fetch", fetchMock);

    render(<PackCalculator />);
    await user.type(await screen.findByLabelText("Add pack size"), "750");
    await user.click(screen.getByRole("button", { name: "Add" }));

    expect(await screen.findByText("750")).toBeInTheDocument();
    expect(fetchMock).toHaveBeenLastCalledWith("/api/packs", expect.objectContaining({ body: '{"sizes":[250,500,750]}' }));

    await user.click(screen.getByRole("button", { name: "Remove 250" }));

    expect(fetchMock).toHaveBeenLastCalledWith("/api/packs", expect.objectContaining({ body: '{"sizes":[500,750]}' }));
    await waitFor(() => expect(screen.queryByRole("button", { name: "Remove 250" })).not.toBeInTheDocument());
  });

  it("resets pack sizes to the defaults", async () => {
    const user = userEvent.setup();
    const fetchMock = vi.fn((path: string) =>
      Promise.resolve(
        jsonResponse(
          path === "/api/packs/reset"
            ? [{ id: 1, size: 250 }, { id: 2, size: 500 }, { id: 3, size: 1000 }, { id: 4, size: 2000 }, { id: 5, size: 5000 }]
            : [{ id: 1, size: 23 }],
        ),
      ),
    );
    vi.stubGlobal("fetch", fetchMock);

    render(<PackCalculator />);
    await screen.findByText("23");
    await user.click(screen.getByRole("button", { name: "Reset to defaults" }));

    expect(await screen.findByText("5,000")).toBeInTheDocument();
    expect(screen.queryByText("23")).not.toBeInTheDocument();
    expect(fetchMock).toHaveBeenLastCalledWith("/api/packs/reset", { method: "POST" });
  });
});
