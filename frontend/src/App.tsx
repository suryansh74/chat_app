import { useEffect } from "react"
import { BrowserRouter, Routes, Route, Navigate, useNavigate } from "react-router"
import { ThemeProvider } from "@/contexts/ThemeContext"
import { AuthProvider, useAuth } from "@/contexts/AuthContext"
import { ProtectedRoute, GuestRoute } from "@/components/ProtectedRoute"
import { LoginPage } from "@/pages/LoginPage"
import { RegisterPage } from "@/pages/RegisterPage"
import { VerifyEmailPage } from "@/pages/VerifyEmailPage"
import { VerifyOTPPage } from "@/pages/VerifyOTPPage"
import { ForgotPasswordPage } from "@/pages/ForgotPasswordPage"
import { ResetPasswordPage } from "@/pages/ResetPasswordPage"
import { HomePage } from "@/pages/HomePage"
import { Toaster } from "sonner"
import { WebSocketProvider } from "@/contexts/WebSocketContext"

function InitialRedirect() {
  const { user, isVerified, isLoading } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (isLoading) return

    if (!user) {
      navigate("/login", { replace: true })
    } else if (!isVerified) {
      navigate("/verify-email", { replace: true })
    } else {
      navigate("/home", { replace: true })
    }
  }, [user, isVerified, isLoading, navigate])

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    )
  }

  return null
}

function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<InitialRedirect />} />
      
      <Route
        path="/login"
        element={
          <GuestRoute>
            <LoginPage />
          </GuestRoute>
        }
      />
      <Route
        path="/register"
        element={
          <GuestRoute>
            <RegisterPage />
          </GuestRoute>
        }
      />
      
      <Route
        path="/verify-email"
        element={
          <ProtectedRoute requireVerified={false}>
            <VerifyEmailPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/verify-otp"
        element={
          <ProtectedRoute requireVerified={false}>
            <VerifyOTPPage />
          </ProtectedRoute>
        }
      />
      
      <Route
        path="/forgot-password"
        element={
          <GuestRoute>
            <ForgotPasswordPage />
          </GuestRoute>
        }
      />
      <Route
        path="/reset-password"
        element={
          <GuestRoute>
            <ResetPasswordPage />
          </GuestRoute>
        }
      />
      
      <Route
        path="/home"
        element={
          <ProtectedRoute>
            <HomePage />
          </ProtectedRoute>
        }
      />
      
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

function App() {
  return (
    <BrowserRouter>
      <ThemeProvider>
        <AuthProvider>
          <WebSocketProvider>
            <AppRoutes />
            <Toaster position="top-right" richColors />
          </WebSocketProvider>
        </AuthProvider>
      </ThemeProvider>
    </BrowserRouter>
  )
}

export default App