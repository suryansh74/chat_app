import { useState } from "react"
import { useFormContext } from "react-hook-form"
import { cn } from "@/lib/utils"
import { Eye, EyeOff } from "lucide-react"

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
  const [showPassword, setShowPassword] = useState(false)
  const {
    register,
    watch,
    formState: { errors, touchedFields },
  } = useFormContext()

  const fieldValue = watch(name)
  const hasValue = fieldValue !== undefined && fieldValue !== ""
  const error = errors[name]?.message as string | undefined
  const isTouched = touchedFields[name as keyof typeof touchedFields]
  const hasError = isTouched && error
  const isValid = isTouched && hasValue && !error && type !== "checkbox"
  const isPassword = type === "password"

  return (
    <div className="space-y-2">
      <label
        htmlFor={name}
        className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
      >
        {label}
      </label>
      <div className="relative">
        <input
          id={name}
          type={isPassword && showPassword ? "text" : type}
          placeholder={placeholder}
          disabled={disabled}
          {...register(name)}
          className={cn(
            "flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm pr-10",
            hasError && "border-destructive focus-visible:ring-destructive",
            isValid && "border-green-500 focus-visible:ring-green-500",
            !hasError && !isValid && "border-input"
          )}
        />
        {isPassword && (
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
          >
            {showPassword ? (
              <EyeOff className="h-4 w-4" />
            ) : (
              <Eye className="h-4 w-4" />
            )}
          </button>
        )}
      </div>
      {hasError && (
        <p className="text-sm text-destructive">{error}</p>
      )}
    </div>
  )
}