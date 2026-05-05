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
import { MessageSquare } from "lucide-react";

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
    <div className="min-h-screen flex items-center justify-center px-4 py-12 bg-gradient-to-br from-background via-background to-muted/30">
      <div className="w-full max-w-md">
        <div className="flex items-center justify-center gap-2 mb-8">
          <div className="p-2.5 rounded-xl bg-primary/10">
            <MessageSquare className="h-8 w-8 text-primary" />
          </div>
          <span className="text-2xl font-bold tracking-tight">ChatApp</span>
        </div>
        <Card className="border-border/60 shadow-xl shadow-black/5">
          <CardHeader className="space-y-1 pb-6">
            <CardTitle className="text-2xl font-bold">
              Create an account
            </CardTitle>
            <CardDescription>Enter your details to get started</CardDescription>
          </CardHeader>
          <FormProvider {...methods}>
            <form onSubmit={handleSubmit(onSubmit)}>
              <CardContent className="space-y-5">
                {errors.root && (
                  <div className="rounded-lg bg-red-500/10 border border-red-500/20 p-3.5 text-sm text-red-500 flex items-center gap-2">
                    {errors.root.message}
                  </div>
                )}
                <FormField
                  name="name"
                  label="Full Name"
                  placeholder="John Doe"
                />
                <FormField
                  name="email"
                  label="Email"
                  type="email"
                  placeholder="name@example.com"
                />
                <FormField
                  name="password"
                  label="Password"
                  type="password"
                  placeholder="Create a strong password"
                />
                <FormField
                  name="password_confirmation"
                  label="Confirm Password"
                  type="password"
                  placeholder="Confirm your password"
                />
              </CardContent>
              <CardFooter className="flex flex-col gap-4 pt-2">
                <Button
                  type="submit"
                  className="w-full h-11 text-base font-semibold"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? (
                    <span className="flex items-center gap-2">
                      <span className="h-4 w-4 border-2 border-current border-t-transparent rounded-full animate-spin" />
                      Creating account...
                    </span>
                  ) : (
                    "Create account"
                  )}
                </Button>
                <p className="text-center text-sm text-muted-foreground">
                  Already have an account?{" "}
                  <Link
                    to="/login"
                    className="font-semibold text-primary hover:text-primary/80 transition-colors"
                  >
                    Sign in
                  </Link>
                </p>
              </CardFooter>
            </form>
          </FormProvider>
        </Card>
      </div>
    </div>
  );
}
