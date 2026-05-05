import { useState, useEffect, useCallback } from "react"
import { Check, X, UserPlus, CheckCircle, XCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { notificationApi, friendsApi, type NotificationListItem } from "@/lib/api"
import { toast } from "sonner"

type HandledStatus = "accepted" | "rejected" | null

interface NotificationPanelProps {
  onClose: () => void
  onNotificationHandled?: () => void
}

export function NotificationPanel({ onClose, onNotificationHandled }: NotificationPanelProps) {
  const [notifications, setNotifications] = useState<NotificationListItem[]>([])
  const [loading, setLoading] = useState(true)
  const [handling, setHandling] = useState<string | null>(null)
  const [handled, setHandled] = useState<Map<string, HandledStatus>>(new Map())

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
    if (handling) return
    setHandling(id)

    const { error } = await friendsApi.acceptFriendRequest(id)
    setHandling(null)

    if (error) {
      toast.error(error)
      return
    }

    toast.success("Friend request accepted!")

    setNotifications((prev) => prev.filter((n) => n.id !== id))
    setHandled((prev) => new Map(prev).set(id, "accepted"))
    onNotificationHandled?.()
  }

  const handleReject = async (id: string) => {
    if (handling) return
    setHandling(id)

    const { error } = await friendsApi.rejectFriendRequest(id)
    setHandling(null)

    if (error) {
      toast.error(error)
      return
    }

    toast.info("Friend request rejected")

    setNotifications((prev) => prev.filter((n) => n.id !== id))
    setHandled((prev) => new Map(prev).set(id, "rejected"))
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

  const renderHandledNotification = (status: HandledStatus, id: string) => (
    <li key={id} className={`p-3 ${status === "accepted" ? "bg-emerald-500/10" : "bg-red-500/10"}`}>
      <div className="flex items-start gap-2">
        {status === "accepted" ? (
          <CheckCircle className="mt-1 h-4 w-4 text-emerald-500" />
        ) : (
          <XCircle className="mt-1 h-4 w-4 text-red-500" />
        )}
        <div className="flex-1">
          <p className={`font-medium ${status === "accepted" ? "text-emerald-500" : "text-red-500"}`}>
            {status === "accepted" ? "Friend request accepted" : "Friend request rejected"}
          </p>
        </div>
      </div>
    </li>
  )

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
        ) : notifications.length === 0 && handled.size === 0 ? (
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
                          disabled={handling === notification.id}
                          onClick={() => handleAccept(notification.id)}
                        >
                          <Check className="mr-1 h-3 w-3" />
                          {handling === notification.id ? "..." : "Accept"}
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          disabled={handling === notification.id}
                          onClick={() => handleReject(notification.id)}
                        >
                          <X className="mr-1 h-3 w-3" />
                          {handling === notification.id ? "..." : "Reject"}
                        </Button>
                      </div>
                    )}
                  </div>
                </div>
              </li>
            ))}

            {Array.from(handled).map(([id, status]) =>
              renderHandledNotification(status, id)
            )}
          </ul>
        )}
      </div>
    </div>
  )
}
