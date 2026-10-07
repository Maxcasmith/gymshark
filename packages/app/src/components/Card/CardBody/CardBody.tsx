import type { ReactNode } from "react";

interface CardBodyProps {
  children: ReactNode;
}

export function CardBody(props: CardBodyProps) {
  const { children } = props;

  return <div className="w-full px-6 py-4">{children}</div>;
}
