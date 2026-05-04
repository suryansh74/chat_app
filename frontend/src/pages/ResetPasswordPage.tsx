import { useForm, FormProvider } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useNavigate, Link } from "react-router"
import { authApi } from "@/lib/api"
import { resetPasswordSchema, type ResetPasswordInput } from "@/lib/schemas"
import { Button } from "@/components/ui/button"
import { FormField } from "@/components/FormField"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"

export function ResetPasswordPage() {
  const navigate = useNavigate()

  const methods = useForm<ResetPasswordInput>({
    resolver: zodResolver(resetPasswordSchema),
    mode: "onChange",
  })

  const {
    handleSubmit,
    formState: { errors, isSubmitting },
  } = methods

  const onSubmit = async (data: ResetPasswordInput) => {
    const { error } = await authApi.resetPassword(data.password, data.password_confirmation)

    if (error) {
      return
    }

    navigate("/login")
  }

  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold">Reset password</CardTitle>
          <CardDescription>Enter your new password</CardDescription>
        </CardHeader>
        <FormProvider {...methods}>
          <form onSubmit={handleSubmit(onSubmit)}>
            <CardContent className="space-y-4">
              {errors.root && (
                <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                  {errors.root.message}
                </div>
              )}
              <FormField name="password" label="New Password" type="password" />
              <FormField name="password_confirmation" label="Confirm Password" type="password" />
            </CardContent>
            <CardFooter className="flex flex-col space-y-4">
              <Button type="submit" className="w-full" disabled={isSubmitting}>
                {isSubmitting ? "Resetting..." : "Reset Password"}
              </Button>
              <div className="text-center text-sm">
                <span className="text-muted-foreground">Back to </span>
                <Link to="/login" className="text-primary hover:underline">
                  Login
                </Link>
              </div>
            </CardFooter>
          </form>
        </FormProvider>
      </Card>
    </div>
  )
}