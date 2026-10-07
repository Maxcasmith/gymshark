import type { ReactNode } from "react";

interface PlatformLayoutProps {
  title: string;
  children: ReactNode;
}

export function PlatformLayout(props: PlatformLayoutProps) {
  const { title, children } = props;

  return (
    <div className="min-h-screen bg-white text-black antialiased">
      <div className="bg-black px-4 py-2 text-center text-xs font-bold uppercase tracking-widest text-white">
        Whole packs only · fewest items · fewest packs
      </div>
      <header className="border-b border-neutral-200 bg-white px-4 py-5">
        <h1 className="mx-auto max-w-3xl text-2xl font-black uppercase tracking-tight">{title}</h1>
      </header>
      <main className="mx-auto max-w-3xl space-y-6 px-4 py-10">{children}</main>
    </div>
  );
}
