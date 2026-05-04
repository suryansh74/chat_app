import { useAuth } from "@/contexts/AuthContext"
import { useTheme } from "@/contexts/ThemeContext"
import { Button } from "@/components/ui/button"
import { Moon, Sun } from "lucide-react"

export function Navbar() {
  const { user, logout } = useAuth()
  const { resolvedTheme, setTheme } = useTheme()

  const handleLogout = async () => {
    console.log("[Navbar] Calling logout...")
    await logout()
    // Clear cookie manually as fallback
    document.cookie = "session_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT"
    console.log("[Navbar] Logout complete, cookie cleared")
    window.location.href = "/login"
  }

  const toggleTheme = () => {
    setTheme(resolvedTheme === "dark" ? "light" : "dark")
  }

  return (
    <nav className="flex h-16 items-center justify-between border-b px-6">
      <div className="flex items-center gap-2">
        <span className="text-xl font-bold">ChatApp</span>
      </div>

      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={toggleTheme}>
          {resolvedTheme === "dark" ? (
            <Sun className="h-5 w-5" />
          ) : (
            <Moon className="h-5 w-5" />
          )}
        </Button>

        <div className="flex items-center gap-2">
          <span className="text-sm text-muted-foreground">{user?.email}</span>
          <Button variant="outline" size="sm" onClick={handleLogout}>
            Logout
          </Button>
        </div>
      </div>
    </nav>
  )
}