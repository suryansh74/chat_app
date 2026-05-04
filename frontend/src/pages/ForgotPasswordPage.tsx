import { useForm, FormProvider } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useNavigate, Link } from "react-router"
import { authApi } from "@/lib/api"
import { forgotPasswordSchema, type ForgotPasswordInput } from "@/lib/schemas"
import { Button } from "@/components/ui/button"
import { FormField } from "@/components/FormField"
import { toast } from "sonner"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"

export function ForgotPasswordPage() {
  const navigate = useNavigate()

  const methods = useForm<ForgotPasswordInput>({
    resolver: zodResolver(forgotPasswordSchema),
    mode: "onChange",
  })

  const {
    handleSubmit,
    formState: { errors, isSubmitting },
  } = methods

  const onSubmit = async (data: ForgotPasswordInput) => {
    const { error } = await authApi.forgotPassword(data.email)

    if (error) {
      toast.error(error)
      return
    }

    toast.success("OTP sent to your email!")
    navigate(`/verify-otp?mode=password-reset&email=${encodeURIComponent(data.email)}`)
  }

  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold">Forgot password</CardTitle>
          <CardDescription>
            Enter your email and we&apos;ll send you an OTP to reset your password
          </CardDescription>
        </CardHeader>
        <FormProvider {...methods}>
          <form onSubmit={handleSubmit(onSubmit)}>
            <CardContent className="space-y-4">
              {errors.root && (
                <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                  {errors.root.message}
                </div>
              )}
              <FormField name="email" label="Email" type="email" placeholder="m@example.com" />
            </CardContent>
            <CardFooter className="flex flex-col space-y-4">
              <Button type="submit" className="w-full" disabled={isSubmitting}>
                {isSubmitting ? "Sending OTP..." : "Send OTP"}
              </Button>
              <div className="text-center text-sm">
                <span className="text-muted-foreground">Remember your password? </span>
                <Link to="/login" className="text-primary hover:underline">
                  Sign in
                </Link>
              </div>
            </CardFooter>
          </form>
        </FormProvider>
      </Card>
    </div>
  )
}