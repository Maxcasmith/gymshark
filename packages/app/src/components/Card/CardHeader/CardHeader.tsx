import type { ReactNode } from "react";

interface CardHeaderProps {
  children: ReactNode;
}

export function CardHeader(props: CardHeaderProps) {
  const { children } = props;

  return <div className="flex w-full justify-between px-6 pt-6 pb-4">{children}</div>;
}
