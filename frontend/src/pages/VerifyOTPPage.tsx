import { useState, useRef, useEffect } from "react"
import { useNavigate, useSearchParams } from "react-router"
import { useAuth } from "@/contexts/AuthContext"
import { authApi } from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { toast } from "sonner"

type Mode = "email-verification" | "password-reset"

export function VerifyOTPPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const mode = (searchParams.get("mode") as Mode) || "email-verification"
  const email = searchParams.get("email") || ""
  
  const { checkAuth } = useAuth()
  const [otp, setOtp] = useState(["", "", "", "", "", ""])
  const inputRefs = useRef<(HTMLInputElement | null)[]>([])
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    if (inputRefs.current[0]) {
      inputRefs.current[0].focus()
    }
  }, [])

  const handlePaste = (e: React.ClipboardEvent) => {
    e.preventDefault()
    const pasteData = e.clipboardData.getData("text")
    const digits = pasteData.replace(/\D/g, "").slice(0, 6).split("")
    
    if (digits.length === 6) {
      setOtp(digits)
      inputRefs.current[5]?.focus()
    }
  }

  const handleChange = (index: number, value: string) => {
    if (!/^\d*$/.test(value)) return

    const newOtp = [...otp]
    newOtp[index] = value
    setOtp(newOtp)

    if (value && index < 5) {
      inputRefs.current[index + 1]?.focus()
    }
  }

  const handleKeyDown = (index: number, e: React.KeyboardEvent) => {
    if (e.key === "Backspace" && !otp[index] && index > 0) {
      inputRefs.current[index - 1]?.focus()
    }
  }

  const handleVerify = async () => {
    const otpValue = otp.join("")
    if (otpValue.length !== 6) {
      toast.error("Please enter all 6 digits")
      return
    }

    setIsLoading(true)

    if (mode === "password-reset") {
      const { error } = await authApi.verifyResetOTP(email, otpValue)
      if (error) {
        toast.error(error)
        setIsLoading(false)
        return
      }
      navigate(`/reset-password?email=${encodeURIComponent(email)}`)
      return
    }

    const { error } = await authApi.verifyOTP(otpValue)
    if (error) {
      toast.error(error)
      setIsLoading(false)
      return
    }

    toast.success("Email verified successfully!")
    setIsLoading(false)
    
    await checkAuth()
    navigate("/home")
  }

  const isPasswordReset = mode === "password-reset"

  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold">
            {isPasswordReset ? "Reset Password" : "Verify OTP"}
          </CardTitle>
          <CardDescription>
            {isPasswordReset
              ? "Enter the OTP sent to your email to reset your password"
              : "Enter the 6-digit code sent to your email"}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex justify-center gap-2" onPaste={handlePaste}>
            {otp.map((digit, index) => (
              <Input
                key={index}
                ref={(el: HTMLInputElement | null) => { inputRefs.current[index] = el }}
                type="text"
                inputMode="numeric"
                maxLength={1}
                value={digit}
                onChange={(e) => handleChange(index, e.target.value)}
                onKeyDown={(e) => handleKeyDown(index, e)}
                className="h-12 w-12 text-center text-lg"
              />
            ))}
          </div>

          <Button onClick={handleVerify} className="w-full" disabled={isLoading}>
            {isLoading ? "Verifying..." : isPasswordReset ? "Verify & Continue" : "Verify"}
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}