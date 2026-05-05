import { useState, useEffect } from "react"
import { Bell, Search, UserPlus, MessageCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { notificationApi } from "@/lib/api"

interface HeaderProps {
  onOpenFriendRequest: () => void
  onOpenNotifications: () => void
  onOpenSearch: () => void
}

export function Header({ onOpenFriendRequest, onOpenNotifications, onOpenSearch }: HeaderProps) {
  const [unreadCount, setUnreadCount] = useState(0)

  const loadUnreadCount = async () => {
    const { data } = await notificationApi.getUnreadCount()
    if (data) {
      setUnreadCount(data.unread_count)
    }
  }

  useEffect(() => {
    loadUnreadCount()
    const interval = setInterval(loadUnreadCount, 30000)
    return () => clearInterval(interval)
  }, [])

  return (
    <header className="flex items-center justify-between border-b bg-background px-4 py-3">
      <div className="flex items-center gap-2">
        <MessageCircle className="h-6 w-6" />
        <h1 className="text-xl font-semibold">Chat</h1>
      </div>
      
      <div className="flex items-center gap-2">
        <Button variant="ghost" size="icon" onClick={onOpenSearch}>
          <Search className="h-5 w-5" />
        </Button>
        
        <Button variant="ghost" size="icon" onClick={onOpenFriendRequest}>
          <UserPlus className="h-5 w-5" />
        </Button>
        
        <Button variant="ghost" size="icon" onClick={onOpenNotifications} className="relative">
          <Bell className="h-5 w-5" />
          {unreadCount > 0 && (
            <span className="absolute -right-1 -top-1 flex h-5 w-5 items-center justify-center rounded-full bg-destructive text-xs text-destructive-foreground">
              {unreadCount > 9 ? "9+" : unreadCount}
            </span>
          )}
        </Button>
      </div>
    </header>
  )
}