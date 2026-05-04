import { useForm, FormProvider } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate, Link } from "react-router";
import { useAuth } from "@/contexts/AuthContext";
import { registerSchema, type RegisterInput } from "@/lib/schemas";
import { Button } from "@/components/ui/button";
import { FormField } from "@/components/FormField";
import { toast } from "sonner";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export function RegisterPage() {
  const navigate = useNavigate();
  const { register: registerUser } = useAuth();

  const methods = useForm<RegisterInput>({
    resolver: zodResolver(registerSchema),
    mode: "onChange",
    defaultValues: {
      name: "Suryansh Awasthi",
      email: "suryanshawasthi56@gmail.com",
      password: "Krazymon123#",
      password_confirmation: "Krazymon123#",
    },
  });

  const {
    handleSubmit,
    formState: { errors, isSubmitting },
  } = methods;

  const onSubmit = async (data: RegisterInput) => {
    const result = await registerUser(
      data.name,
      data.email,
      data.password,
      data.password_confirmation,
    );

    if (result.error) {
      toast.error(result.error);
      return;
    }

    toast.success("Account created! Please verify your email.");
    navigate("/verify-email");
  };

  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold">
            Create an account
          </CardTitle>
          <CardDescription>Enter your details to get started</CardDescription>
        </CardHeader>
        <FormProvider {...methods}>
          <form onSubmit={handleSubmit(onSubmit)}>
            <CardContent className="space-y-4">
              {errors.root && (
                <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                  {errors.root.message}
                </div>
              )}
              <FormField name="name" label="Name" placeholder="John Doe" />
              <FormField
                name="email"
                label="Email"
                type="email"
                placeholder="m@example.com"
              />
              <FormField name="password" label="Password" type="password" />
              <FormField
                name="password_confirmation"
                label="Confirm Password"
                type="password"
              />
            </CardContent>
            <CardFooter className="flex flex-col space-y-4">
              <Button type="submit" className="w-full" disabled={isSubmitting}>
                {isSubmitting ? "Creating account..." : "Create account"}
              </Button>
              <div className="text-center text-sm">
                <span className="text-muted-foreground">
                  Already have an account?{" "}
                </span>
                <Link to="/login" className="text-primary hover:underline">
                  Sign in
                </Link>
              </div>
            </CardFooter>
          </form>
        </FormProvider>
      </Card>
    </div>
  );
}

