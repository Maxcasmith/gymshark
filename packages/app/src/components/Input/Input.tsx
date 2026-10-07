import type { InputHTMLAttributes } from "react";
import { cn } from "@/lib/utils";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  errors?: string[];
}

export function Input(props: InputProps) {
  const { name, label = name, errors, className, ...rest } = props;
  const hasErrors = Boolean(errors?.length);

  return (
    <div className="flex flex-col">
      <label htmlFor={name} className="block text-xs font-bold uppercase tracking-wider text-black">
        {label}
        <input
          id={name}
          name={name}
          aria-invalid={hasErrors}
          className={cn(
            "mt-2 block h-12 w-full rounded-md border border-neutral-300 bg-white px-4 text-base font-normal normal-case tracking-normal text-black placeholder:text-neutral-400 focus:border-black focus:outline-none disabled:bg-neutral-100 disabled:text-neutral-500",
            hasErrors && "border-red-600 bg-red-50 focus:border-red-600",
            className,
          )}
          {...rest}
        />
      </label>
      {hasErrors && (
        <div className="mt-2 text-xs text-red-600">
          <span>Please address the following errors</span>
          <ul className="m-0 list-disc pl-4">
            {errors?.map((error) => (
              <li key={error}>{error}</li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}
