import type { ReactNode } from "react";

interface TagProps {
  children: ReactNode;
}

export function Tag(props: TagProps) {
  const { children } = props;

  return (
    <span className="inline-block rounded-sm bg-neutral-100 px-3 py-1 text-xs font-bold uppercase tracking-wider text-black">
      {children}
    </span>
  );
}
