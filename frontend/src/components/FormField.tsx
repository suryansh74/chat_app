import { useFormContext } from "react-hook-form"
import { cn } from "@/lib/utils"

interface FormFieldProps {
  name: string
  label: string
  type?: "text" | "email" | "password" | "checkbox"
  placeholder?: string
  disabled?: boolean
}

export function FormField({
  name,
  label,
  type = "text",
  placeholder,
  disabled,
}: FormFieldProps) {
  const {
    register,
    formState: { errors, touchedFields },
  } = useFormContext()

  const error = errors[name]?.message as string | undefined
  const isTouched = touchedFields[name as keyof typeof touchedFields]
  const hasError = isTouched && error
  const isValid = isTouched && !error && type !== "checkbox"

  return (
    <div className="space-y-2">
      <label
        htmlFor={name}
        className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
      >
        {label}
      </label>
      <input
        id={name}
        type={type}
        placeholder={placeholder}
        disabled={disabled}
        {...register(name)}
        className={cn(
          "flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
          hasError && "border-destructive focus-visible:ring-destructive",
          isValid && "border-green-500 focus-visible:ring-green-500",
          !hasError && !isValid && "border-input"
        )}
      />
      {hasError && (
        <p className="text-sm text-destructive">{error}</p>
      )}
    </div>
  )
}