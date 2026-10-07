import { useEffect, useState, type FormEvent } from "react";
import { Button } from "@/components/Button/Button";
import { Card, CardBody, CardFooter, CardHeader } from "@/components/Card";
import { DataTable, type DataTableColumn } from "@/components/DataTable/DataTable";
import { Input } from "@/components/Input/Input";
import { Spinner } from "@/components/Spinner/Spinner";
import { Tag } from "@/components/Tag/Tag";

interface Pack {
  id: number;
  size: number;
}

interface OrderPack {
  packSize: number;
  quantity: number;
}

interface Order {
  id: number;
  requestedItems: number;
  shippedItems: number;
  excessItems: number;
  totalPacks: number;
  packs: OrderPack[];
}

interface ErrorResponse {
  error: string;
}

const formatNumber = (value: number) => value.toLocaleString();

const packColumns: DataTableColumn<OrderPack>[] = [
  { key: "packSize", header: "Pack size", cell: (orderPack) => formatNumber(orderPack.packSize) },
  { key: "quantity", header: "Quantity", cell: (orderPack) => formatNumber(orderPack.quantity), align: "right" },
];

function isErrorResponse(body: unknown): body is ErrorResponse {
  return typeof body === "object" && body !== null && "error" in body && typeof body.error === "string";
}

async function request<TResponse>(path: string, init?: RequestInit): Promise<TResponse> {
  const response = await fetch(path, init);
  const body: unknown = await response.json();
  if (!response.ok) {
    throw new Error(isErrorResponse(body) ? body.error : response.statusText);
  }
  return body as TResponse;
}

export function PackCalculator() {
  const [packs, setPacks] = useState<Pack[] | null>(null);
  const [requestedItems, setRequestedItems] = useState("");
  const [order, setOrder] = useState<Order | null>(null);
  const [errorMessage, setErrorMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [newPackSize, setNewPackSize] = useState("");
  const [packErrorMessage, setPackErrorMessage] = useState("");
  const [isSavingPacks, setIsSavingPacks] = useState(false);

  useEffect(() => {
    request<Pack[]>("/api/packs").then(setPacks, (error: Error) => setPackErrorMessage(error.message));
  }, []);

  async function updatePacks(path: string, init: RequestInit) {
    setIsSavingPacks(true);
    setPackErrorMessage("");
    try {
      setPacks(await request<Pack[]>(path, init));
      return true;
    } catch (error) {
      setPackErrorMessage(error instanceof Error ? error.message : String(error));
      return false;
    } finally {
      setIsSavingPacks(false);
    }
  }

  function savePackSizes(sizes: number[]) {
    return updatePacks("/api/packs", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ sizes }),
    });
  }

  async function handleAddPack(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (await savePackSizes([...(packs ?? []).map((pack) => pack.size), Number(newPackSize)])) {
      setNewPackSize("");
    }
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsSubmitting(true);
    setErrorMessage("");
    try {
      setOrder(
        await request<Order>("/api/orders", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ requestedItems: Number(requestedItems) }),
        }),
      );
    } catch (error) {
      setOrder(null);
      setErrorMessage(error instanceof Error ? error.message : String(error));
    } finally {
      setIsSubmitting(false);
    }
  }

  const stats = order
    ? [
        { label: "Requested", value: order.requestedItems },
        { label: "Shipped", value: order.shippedItems },
        { label: "Excess", value: order.excessItems },
      ]
    : [];

  return (
    <>
      <Card>
        <form onSubmit={handleAddPack}>
          <CardHeader>
            <div className="space-y-3">
              <h2 className="text-xl font-black uppercase tracking-tight">Pack sizes</h2>
              {packs ? (
                <div className="flex flex-wrap gap-2">
                  {packs.map((pack) => (
                    <Tag key={pack.id}>
                      {formatNumber(pack.size)}
                      <button
                        type="button"
                        aria-label={`Remove ${formatNumber(pack.size)}`}
                        disabled={isSavingPacks || packs.length === 1}
                        onClick={() => savePackSizes(packs.filter(({ id }) => id !== pack.id).map(({ size }) => size))}
                        className="ml-2 cursor-pointer text-neutral-500 hover:text-black disabled:cursor-not-allowed disabled:opacity-40"
                      >
                        ×
                      </button>
                    </Tag>
                  ))}
                </div>
              ) : (
                !packErrorMessage && <Spinner />
              )}
            </div>
          </CardHeader>
          <CardBody>
            <Input
              name="newPackSize"
              label="Add pack size"
              type="number"
              min={1}
              required
              placeholder="e.g. 750"
              value={newPackSize}
              onChange={(event) => setNewPackSize(event.target.value)}
            />
            {packErrorMessage && (
              <p role="alert" className="mt-3 text-sm font-semibold text-red-600">
                {packErrorMessage}
              </p>
            )}
          </CardBody>
          <CardFooter>
            <span className="text-xs text-neutral-500">Changes apply to the next order</span>
            <div className="flex gap-3">
              <Button
                type="button"
                disabled={isSavingPacks || !packs}
                onClick={() => updatePacks("/api/packs/reset", { method: "POST" })}
                className="border border-black bg-white text-black hover:bg-neutral-100"
              >
                Reset to defaults
              </Button>
              <Button type="submit" disabled={isSavingPacks || !packs}>
                {isSavingPacks ? "Saving…" : "Add"}
              </Button>
            </div>
          </CardFooter>
        </form>
      </Card>

      <Card>
        <form onSubmit={handleSubmit}>
          <CardHeader>
            <h2 className="text-xl font-black uppercase tracking-tight">Calculate packs</h2>
          </CardHeader>
          <CardBody>
            <Input
              name="requestedItems"
              label="Items ordered"
              type="number"
              min={1}
              required
              placeholder="e.g. 12001"
              value={requestedItems}
              onChange={(event) => setRequestedItems(event.target.value)}
            />
            {errorMessage && (
              <p role="alert" className="mt-3 text-sm font-semibold text-red-600">
                {errorMessage}
              </p>
            )}
          </CardBody>
          <CardFooter>
            <span className="text-xs text-neutral-500">Fewest items first, then fewest packs</span>
            <Button type="submit" disabled={isSubmitting || !packs}>
              {isSubmitting ? "Calculating…" : "Calculate"}
            </Button>
          </CardFooter>
        </form>
      </Card>

      {order && (
        <Card>
          <CardHeader>
            <dl className="grid w-full grid-cols-3 gap-3 text-center">
              {stats.map((stat) => (
                <div key={stat.label}>
                  <dt className="text-xs font-bold uppercase tracking-wider text-neutral-500">{stat.label}</dt>
                  <dd className="mt-1 text-3xl font-black">{formatNumber(stat.value)}</dd>
                </div>
              ))}
            </dl>
          </CardHeader>
          <CardBody>
            <DataTable
              caption="Packs to ship"
              columns={packColumns}
              rows={order.packs}
              rowKey={(orderPack) => orderPack.packSize}
            />
          </CardBody>
          <CardFooter>
            <span className="text-xs font-bold uppercase tracking-wider">Total packs</span>
            <span className="text-xl font-black">{formatNumber(order.totalPacks)}</span>
          </CardFooter>
        </Card>
      )}
    </>
  );
}
