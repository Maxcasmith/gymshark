import type { ReactNode } from "react";

interface CardFooterProps {
  children: ReactNode;
}

export function CardFooter(props: CardFooterProps) {
  const { children } = props;

  return (
    <div className="flex w-full flex-wrap items-center justify-between gap-4 border-t border-neutral-200 bg-neutral-50 px-6 py-4">
      {children}
    </div>
  );
}
