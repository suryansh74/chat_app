import { useState } from "react"
import { useFormContext } from "react-hook-form"
import { cn } from "@/lib/utils"
import { Eye, EyeOff, AlertCircle, CheckCircle2 } from "lucide-react"

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
    <div className="space-y-1.5">
      <label
        htmlFor={name}
        className="text-sm font-semibold text-foreground/90 peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
      >
        {label}
      </label>
      <div className="relative group">
        <input
          id={name}
          type={isPassword && showPassword ? "text" : type}
          placeholder={placeholder}
          disabled={disabled}
          {...register(name)}
          className={cn(
            "flex h-11 w-full rounded-lg border bg-background px-4 py-2 text-sm ring-offset-background transition-all duration-200",
            "placeholder:text-muted-foreground/60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2",
            "disabled:cursor-not-allowed disabled:opacity-50",
            isPassword && "pr-12",
            !isPassword && "pr-4",
            hasError && "border-red-500 focus-visible:ring-red-500 focus-visible:border-red-500",
            isValid && "border-emerald-500 focus-visible:ring-emerald-500 focus-visible:border-emerald-500",
            !hasError && !isValid && "border-border hover:border-muted-foreground/30 focus-visible:ring-primary"
          )}
        />
        {isPassword && (
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors p-1 rounded-md hover:bg-muted"
          >
            {showPassword ? (
              <EyeOff className="h-4 w-4" />
            ) : (
              <Eye className="h-4 w-4" />
            )}
          </button>
        )}
        {!isPassword && hasError && (
          <div className="absolute right-3 top-1/2 -translate-y-1/2">
            <AlertCircle className="h-4 w-4 text-red-500" />
          </div>
        )}
        {!isPassword && isValid && (
          <div className="absolute right-3 top-1/2 -translate-y-1/2">
            <CheckCircle2 className="h-4 w-4 text-emerald-500" />
          </div>
        )}
      </div>
      {hasError && (
        <p className="text-sm font-medium text-red-500 flex items-center gap-1.5 animate-in fade-in slide-in-from-top-1">
          <AlertCircle className="h-3.5 w-3.5" />
          {error}
        </p>
      )}
    </div>
  )
}