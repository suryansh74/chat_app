import { useAuth } from "@/contexts/AuthContext"
import { useTheme } from "@/contexts/ThemeContext"
import { Button } from "@/components/ui/button"
import { Moon, Sun, UserPlus, Bell } from "lucide-react"
import { toast } from "sonner"
import { notificationApi, createWebSocket } from "@/lib/api"
import { useState, useEffect, useRef } from "react"

interface NavbarProps {
  onOpenFriendRequest?: () => void
  onOpenNotifications?: () => void
  onNotificationCountChange?: (count: number) => void
  onNotificationHandled?: () => void
}

export function Navbar({ onOpenFriendRequest, onOpenNotifications, onNotificationCountChange, onNotificationHandled }: NavbarProps) {
  const { user, logout } = useAuth()
  const { resolvedTheme, setTheme } = useTheme()
  const [unreadCount, setUnreadCount] = useState(0)
  const wsRef = useRef<WebSocket | null>(null)
  const initializedRef = useRef(false)

  const loadUnreadCount = async () => {
    if (!user?.id) return
    const { data } = await notificationApi.getUnreadCount()
    if (data) {
      setUnreadCount(data.unread_count)
      onNotificationCountChange?.(data.unread_count)
    }
  }

  useEffect(() => {
    if (!user?.id || initializedRef.current) return

    initializedRef.current = true

    // Load initial count
    loadUnreadCount()

    // Connect to WebSocket for real-time notifications
    const ws = createWebSocket(user.id)
    wsRef.current = ws

    ws.onopen = () => {
      console.log("[Navbar] WebSocket connected")
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        console.log("[Navbar] WebSocket message:", data)
        
        if (data.type === "new_notification") {
          toast.info(data.data.content)
          setUnreadCount(prev => prev + 1)
          onNotificationCountChange?.(unreadCount + 1)
        } else if (data.type === "new_friend") {
          toast.success(data.data.content)
          loadUnreadCount()
        }
      } catch (e) {
        console.error("[Navbar] Failed to parse WS message:", e)
      }
    }

    ws.onerror = (error) => {
      console.error("[Navbar] WebSocket error:", error)
    }

    ws.onclose = () => {
      console.log("[Navbar] WebSocket closed")
      wsRef.current = null
    }

    return () => {
      initializedRef.current = false
      if (wsRef.current) {
        wsRef.current.close()
        wsRef.current = null
      }
    }
  }, [user?.id])

  const handleNotificationsClick = async () => {
    onOpenNotifications?.()
    loadUnreadCount()
    onNotificationHandled?.()
  }

  const handleLogout = async () => {
    console.log("[Navbar] Calling logout...")
    
    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }
    
    const result = await logout()
    
    if (result.error) {
      toast.error(result.error)
      return
    }

    if (result.success) {
      toast.success(result.success)
    }
    
    console.log("[Navbar] Logout complete")
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
        {onOpenFriendRequest && (
          <Button variant="ghost" size="icon" onClick={onOpenFriendRequest} title="Add Friend">
            <UserPlus className="h-5 w-5" />
          </Button>
        )}

        {onOpenNotifications && (
          <Button 
            variant="ghost" 
            size="icon" 
            onClick={handleNotificationsClick} 
            title="Notifications" 
            className="relative"
          >
            <Bell className="h-5 w-5" />
            {unreadCount > 0 && (
              <span className="absolute -right-1 -top-1 flex h-4 w-4 items-center justify-center rounded-full bg-destructive text-[10px] text-destructive-foreground">
                {unreadCount > 9 ? "9+" : unreadCount}
              </span>
            )}
          </Button>
        )}

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