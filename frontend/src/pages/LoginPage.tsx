import { useForm, FormProvider } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate, Link } from "react-router";
import { useAuth } from "@/contexts/AuthContext";
import { loginSchema, type LoginInput } from "@/lib/schemas";
import { Button } from "@/components/ui/button";
import { FormField } from "@/components/FormField";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export function LoginPage() {
  const navigate = useNavigate();
  const { login } = useAuth();

  const methods = useForm<LoginInput>({
    resolver: zodResolver(loginSchema),
    mode: "onChange",
    defaultValues: {
      email: "suryanshawasthi56@gmail.com",
      password: "Krazymon123#",
      remember_me: false,
    },
  });

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = methods;

  const onSubmit = async (data: LoginInput) => {
    const result = await login(data.email, data.password, data.remember_me);

    if (result.error) {
      return;
    }

    navigate("/home");
  };

  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold">Sign in</CardTitle>
          <CardDescription>
            Enter your email and password to access your account
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
              <FormField
                name="email"
                label="Email"
                type="email"
                placeholder="m@example.com"
              />
              <FormField name="password" label="Password" type="password" />
              <div className="flex items-center space-x-2">
                <input
                  id="remember_me"
                  type="checkbox"
                  {...register("remember_me")}
                  className="h-4 w-4 rounded border-gray-300"
                />
                <label htmlFor="remember_me" className="text-sm font-normal">
                  Remember me for 30 days
                </label>
              </div>
            </CardContent>
            <CardFooter className="flex flex-col space-y-4">
              <Button type="submit" className="w-full" disabled={isSubmitting}>
                {isSubmitting ? "Signing in..." : "Sign in"}
              </Button>
              <div className="text-center text-sm">
                <span className="text-muted-foreground">
                  Forgot your password?{" "}
                </span>
                <Link
                  to="/forgot-password"
                  className="text-primary hover:underline"
                >
                  Reset
                </Link>
              </div>
              <div className="text-center text-sm">
                <span className="text-muted-foreground">
                  Don&apos;t have an account?{" "}
                </span>
                <Link to="/register" className="text-primary hover:underline">
                  Sign up
                </Link>
              </div>
            </CardFooter>
          </form>
        </FormProvider>
      </Card>
    </div>
  );
}

