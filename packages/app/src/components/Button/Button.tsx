import type { ButtonHTMLAttributes, ReactNode } from "react";
import { cn } from "@/lib/utils";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  children: ReactNode;
}

export function Button(props: ButtonProps) {
  const { children, className, ...rest } = props;

  return (
    <button
      className={cn(
        "inline-flex h-12 cursor-pointer items-center justify-center rounded-none bg-black px-8 text-sm font-bold uppercase tracking-wider text-white transition-colors hover:bg-neutral-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-black disabled:cursor-not-allowed disabled:bg-neutral-300 disabled:text-neutral-500",
        className,
      )}
      {...rest}
    >
      {children}
    </button>
  );
}
