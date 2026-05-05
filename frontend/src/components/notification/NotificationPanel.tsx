import { useState, useEffect, useCallback } from "react"
import { Check, X, UserPlus } from "lucide-react"
import { Button } from "@/components/ui/button"
import { notificationApi, friendsApi, type NotificationListItem } from "@/lib/api"
import { toast } from "sonner"

interface NotificationPanelProps {
  onClose: () => void
  onNotificationHandled?: () => void
}

export function NotificationPanel({ onClose, onNotificationHandled }: NotificationPanelProps) {
  const [notifications, setNotifications] = useState<NotificationListItem[]>([])
  const [loading, setLoading] = useState(true)

  const loadNotifications = useCallback(async () => {
    setLoading(true)
    const { data } = await notificationApi.getNotifications()
    if (data?.notifications) {
      setNotifications(data.notifications)
    }
    setLoading(false)
  }, [])

  useEffect(() => {
    loadNotifications()
  }, [loadNotifications])

  const handleAccept = async (id: string) => {
    const { error } = await friendsApi.acceptFriendRequest(id)
    if (error) {
      toast.error(error)
      return
    }
    
    toast.success("Friend request accepted!")
    
    // Remove from list and refresh
    setNotifications((prev) => prev.filter((n) => n.id !== id))
    loadNotifications()
    
    // Notify parent to refresh friends list
    onNotificationHandled?.()
  }

  const handleReject = async (id: string) => {
    const { error } = await friendsApi.rejectFriendRequest(id)
    if (error) {
      toast.error(error)
      return
    }
    
    toast.info("Friend request rejected")
    
    // Remove from list and refresh
    setNotifications((prev) => prev.filter((n) => n.id !== id))
    loadNotifications()
    
    // Notify parent to update badge count
    onNotificationHandled?.()
  }

  const formatTime = (dateStr: string) => {
    const date = new Date(dateStr)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()
    const diffMins = Math.floor(diffMs / 60000)
    const diffHours = Math.floor(diffMs / 3600000)

    if (diffMins < 1) return "now"
    if (diffMins < 60) return `${diffMins}m ago`
    if (diffHours < 24) return `${diffHours}h ago`
    return date.toLocaleDateString()
  }

  return (
    <div className="flex h-full flex-col border-r">
      <div className="flex items-center justify-between border-b p-3">
        <h2 className="font-semibold">Notifications</h2>
        <Button variant="ghost" size="icon" onClick={onClose}>
          <X className="h-4 w-4" />
        </Button>
      </div>

      <div className="flex-1 overflow-y-auto">
        {loading ? (
          <div className="flex items-center justify-center p-4">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          </div>
        ) : notifications.length === 0 ? (
          <div className="p-4 text-center text-muted-foreground">
            No notifications
          </div>
        ) : (
          <ul className="divide-y">
            {notifications.map((notification) => (
              <li key={notification.id} className={`p-3 ${notification.is_read ? "bg-muted/30" : ""}`}>
                <div className="flex items-start gap-2">
                  <UserPlus className="mt-1 h-4 w-4" />
                  <div className="flex-1">
                    <p>
                      <span className="font-medium">{notification.from_user}</span>{" "}
                      {notification.content}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {formatTime(notification.created_at)}
                    </p>

                    {notification.type === "FRIEND_REQUEST" && (
                      <div className="mt-2 flex gap-2">
                        <Button
                          size="sm"
                          onClick={() => handleAccept(notification.id)}
                        >
                          <Check className="mr-1 h-3 w-3" />
                          Accept
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleReject(notification.id)}
                        >
                          <X className="mr-1 h-3 w-3" />
                          Reject
                        </Button>
                      </div>
                    )}
                  </div>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}