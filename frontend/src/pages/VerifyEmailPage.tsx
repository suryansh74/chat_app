import { useState } from "react"
import { useNavigate } from "react-router"
import { useAuth } from "@/contexts/AuthContext"
import { authApi } from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { toast } from "sonner"

export function VerifyEmailPage() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const [isLoading, setIsLoading] = useState(false)

  const handleSendOTP = async () => {
    setIsLoading(true)

    const { error } = await authApi.sendOTP()

    if (error) {
      toast.error(error)
      setIsLoading(false)
      return
    }

    toast.success("OTP sent to your email!")
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

          <Button onClick={handleSendOTP} className="w-full" disabled={isLoading}>
            {isLoading ? "Sending OTP..." : "Send OTP"}
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}