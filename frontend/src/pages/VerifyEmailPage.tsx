import { useState } from "react"
import { useNavigate } from "react-router"
import { useAuth } from "@/contexts/AuthContext"
import { authApi } from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export function VerifyEmailPage() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const [isLoading, setIsLoading] = useState(false)
  const [message, setMessage] = useState("")
  const [error, setError] = useState("")

  const handleSendOTP = async () => {
    setIsLoading(true)
    setMessage("")
    setError("")

    const { error: apiError } = await authApi.sendOTP()

    if (apiError) {
      setError(apiError)
      setIsLoading(false)
      return
    }

    setMessage("OTP sent to your email!")
    setIsLoading(false)
    navigate("/verify-otp")
  }

  if (!user) {
    return null
  }

  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold">Verify your email</CardTitle>
          <CardDescription>
            We need to verify your email address before you can continue
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="rounded-md bg-muted p-4">
            <p className="text-sm text-muted-foreground">Your email address</p>
            <p className="font-medium">{user.email}</p>
          </div>

          {error && (
            <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
              {error}
            </div>
          )}

          {message && (
            <div className="rounded-md bg-green-500/10 p-3 text-sm text-green-600 dark:text-green-400">
              {message}
            </div>
          )}

          <Button onClick={handleSendOTP} className="w-full" disabled={isLoading}>
            {isLoading ? "Sending OTP..." : "Send OTP"}
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}