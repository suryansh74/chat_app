import { useAuth } from "@/contexts/AuthContext"
import { useTheme } from "@/contexts/ThemeContext"
import { useWebSocket } from "@/contexts/WebSocketContext"
import { Button } from "@/components/ui/button"
import { Moon, Sun, UserPlus, Bell } from "lucide-react"
import { toast } from "sonner"
import { notificationApi } from "@/lib/api"
import { useState, useEffect, useCallback } from "react"

interface NavbarProps {
  notificationRefreshKey?: number
  onOpenFriendRequest?: () => void
  onOpenNotifications?: () => void
  onNotificationCountChange?: (count: number) => void
  onNotificationHandled?: () => void
}

export function Navbar({ notificationRefreshKey, onOpenFriendRequest, onOpenNotifications, onNotificationCountChange }: NavbarProps) {
  const { user, logout } = useAuth()
  const { resolvedTheme, setTheme } = useTheme()
  const { subscribe, unsubscribe } = useWebSocket()
  const [unreadCount, setUnreadCount] = useState(0)

  const loadUnreadCount = useCallback(async () => {
    if (!user?.id) return
    const { data } = await notificationApi.getUnreadCount()
    if (data) {
      setUnreadCount(data.unread_count)
      onNotificationCountChange?.(data.unread_count)
    }
  }, [user?.id, onNotificationCountChange])

  useEffect(() => {
    loadUnreadCount()
  }, [loadUnreadCount])

  useEffect(() => {
    loadUnreadCount()
  }, [notificationRefreshKey])

  useEffect(() => {
    const handleNotification = (data: unknown) => {
      const d = data as { data: { content: string } }
      toast.info(d.data.content)
      setUnreadCount(prev => {
        onNotificationCountChange?.(prev + 1)
        return prev + 1
      })
    }

    subscribe("new_notification", handleNotification)
    return () => unsubscribe("new_notification", handleNotification)
  }, [subscribe, unsubscribe, onNotificationCountChange])

  useEffect(() => {
    const handleNewFriend = () => {
      toast.success("New friend added!")
      loadUnreadCount()
    }

    subscribe("new_friend", handleNewFriend)
    return () => unsubscribe("new_friend", handleNewFriend)
  }, [subscribe, unsubscribe, loadUnreadCount])

  const handleNotificationsClick = async () => {
    setUnreadCount(0)
    await notificationApi.markAllAsRead()
    onOpenNotifications?.()
  }

  const handleLogout = async () => {
    const result = await logout()
    if (result.error) {
      toast.error(result.error)
      return
    }
    if (result.success) {
      toast.success(result.success)
    }
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
