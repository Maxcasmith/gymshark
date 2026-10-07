export function Spinner() {
  return (
    <div role="status" className="flex flex-col items-center justify-center p-8">
      <div className="size-10 animate-spin rounded-full border-4 border-neutral-200 border-t-black" />
      <p className="mt-4 text-xs font-bold uppercase tracking-wider">Loading…</p>
    </div>
  );
}
