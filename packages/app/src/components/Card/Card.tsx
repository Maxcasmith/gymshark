import type { ReactNode } from "react";

interface CardProps {
  children: ReactNode;
}

export function Card(props: CardProps) {
  const { children } = props;

  return <div className="w-full overflow-hidden rounded-md border border-neutral-200 bg-white">{children}</div>;
}
