import { useAuth } from "@/contexts/AuthContext"
import { Navigate } from "react-router"

interface AuthMiddlewareProps {
  children: React.ReactNode
}

interface GuestMiddlewareProps {
  children: React.ReactNode
}

interface VerifiedMiddlewareProps {
  children: React.ReactNode
  redirectTo?: string
}

export function AuthMiddleware({ children }: AuthMiddlewareProps) {
  const { user, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    )
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}

export function GuestMiddleware({ children }: GuestMiddlewareProps) {
  const { user, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    )
  }

  if (user) {
    return <Navigate to="/home" replace />
  }

  return <>{children}</>
}

export function VerifiedMiddleware({ children, redirectTo = "/home" }: VerifiedMiddlewareProps) {
  const { user, isVerified, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    )
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  if (isVerified) {
    return <Navigate to={redirectTo} replace />
  }

  return <>{children}</>
}

export function UnverifiedMiddleware({ children }: VerifiedMiddlewareProps) {
  const { user, isVerified, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    )
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  if (!isVerified) {
    return <Navigate to="/verify-email" replace />
  }

  return <>{children}</>
}