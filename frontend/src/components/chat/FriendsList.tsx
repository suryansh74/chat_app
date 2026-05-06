import { useState, useEffect, useCallback, useRef } from "react"
import { MoreVertical, Trash2 } from "lucide-react"
import { friendsApi, presenceApi, type FriendListItem } from "@/lib/api"
import { useWebSocket } from "@/contexts/WebSocketContext"
import { toast } from "sonner"

interface FriendsListProps {
  onSelectFriend: (friend: FriendListItem) => void
  selectedFriendId?: string
}

export function FriendsList({ onSelectFriend, selectedFriendId }: FriendsListProps) {
  const { subscribe, unsubscribe } = useWebSocket()
  const [friends, setFriends] = useState<FriendListItem[]>([])
  const [loading, setLoading] = useState(true)
  const [onlineMap, setOnlineMap] = useState<Record<string, boolean>>({})
  const [openMenuId, setOpenMenuId] = useState<string | null>(null)
  const [confirmUnfriendId, setConfirmUnfriendId] = useState<string | null>(null)
  const menuRef = useRef<HTMLDivElement>(null)

  const loadFriends = useCallback(async () => {
    setLoading(true)
    const { data } = await friendsApi.getFriends()
    if (data?.friends) {
      setFriends(data.friends)
    }
    setLoading(false)
  }, [])

  useEffect(() => {
    loadFriends()
  }, [loadFriends])

  useEffect(() => {
    if (friends.length === 0) return
    const ids = friends.map((f) => f.friend_id)
    presenceApi.checkOnline(ids).then(({ data }) => {
      if (data?.online) {
        setOnlineMap(data.online)
      }
    })
  }, [friends])

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setOpenMenuId(null)
      }
    }
    document.addEventListener("mousedown", handleClickOutside)
    return () => document.removeEventListener("mousedown", handleClickOutside)
  }, [])

  useEffect(() => {
    const handleNewFriend = () => {
      loadFriends()
    }

    subscribe("new_friend", handleNewFriend)
    return () => unsubscribe("new_friend", handleNewFriend)
  }, [subscribe, unsubscribe, loadFriends])

  useEffect(() => {
    const handleFriendRemoved = (data: unknown) => {
      const d = data as { message: { user_id: string } }
      setFriends((prev) => prev.filter((f) => f.friend_id !== d.message.user_id))
    }

    subscribe("friend_removed", handleFriendRemoved)
    return () => unsubscribe("friend_removed", handleFriendRemoved)
  }, [subscribe, unsubscribe])

  useEffect(() => {
    const handleOnline = (data: unknown) => {
      const d = data as { message: { user_id: string } }
      setOnlineMap((prev) => ({ ...prev, [d.message.user_id]: true }))
    }
    const handleOffline = (data: unknown) => {
      const d = data as { message: { user_id: string } }
      setOnlineMap((prev) => ({ ...prev, [d.message.user_id]: false }))
    }

    subscribe("user_online", handleOnline)
    subscribe("user_offline", handleOffline)
    return () => {
      unsubscribe("user_online", handleOnline)
      unsubscribe("user_offline", handleOffline)
    }
  }, [subscribe, unsubscribe])

  useEffect(() => {
    const handleNewMessage = (data: unknown) => {
      const msgData = data as { message: { from_user_id: string; to_user_id: string; content: string; created_at: string } }
      const msg = msgData.message

      setFriends((prev) => {
        return prev.map((friend) => {
          const isFromThisFriend = friend.friend_id === msg.from_user_id
          const isToThisFriend = friend.friend_id === msg.to_user_id

          if (isFromThisFriend) {
            const shouldIncrement = friend.friend_id !== selectedFriendId
            return {
              ...friend,
              unread_count: shouldIncrement ? friend.unread_count + 1 : friend.unread_count,
              last_message: msg.content,
              last_message_at: msg.created_at,
            }
          }
          if (isToThisFriend) {
            return {
              ...friend,
              unread_count: 0,
              last_message: msg.content,
              last_message_at: msg.created_at,
            }
          }
          return friend
        })
      })
    }

    subscribe("new_message", handleNewMessage)
    return () => unsubscribe("new_message", handleNewMessage)
  }, [subscribe, unsubscribe, selectedFriendId])

  useEffect(() => {
    if (selectedFriendId) {
      setFriends((prev) =>
        prev.map((f) =>
          f.friend_id === selectedFriendId ? { ...f, unread_count: 0 } : f
        )
      )
    }
  }, [selectedFriendId])

  const handleUnfriend = async (friendId: string, friendName: string) => {
    const { error } = await friendsApi.removeFriend(friendId)
    if (error) {
      toast.error(error)
      return
    }

    toast.success(`${friendName} removed from friends`)
    setFriends((prev) => prev.filter((f) => f.friend_id !== friendId))
    setOpenMenuId(null)
    setConfirmUnfriendId(null)
  }

  const formatTime = (dateStr: string) => {
    const date = new Date(dateStr)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()
    const diffMins = Math.floor(diffMs / 60000)
    const diffHours = Math.floor(diffMs / 3600000)
    const diffDays = Math.floor(diffMs / 86400000)

    if (diffMins < 1) return "now"
    if (diffMins < 60) return `${diffMins}m`
    if (diffHours < 24) return `${diffHours}h`
    if (diffDays < 7) return `${diffDays}d`
    return date.toLocaleDateString()
  }

  return (
    <div className="flex h-full flex-col border-r" ref={menuRef}>
      <div className="border-b p-3">
        <h2 className="font-semibold">Friends</h2>
      </div>

      <div className="flex-1 overflow-y-auto">
        {loading ? (
          <div className="flex items-center justify-center p-4">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          </div>
        ) : friends.length === 0 ? (
          <div className="p-4 text-center text-muted-foreground">
            No friends yet. Add some!
          </div>
        ) : (
          <ul className="divide-y">
            {friends.map((friend) => (
              <li key={friend.friend_id} className="relative">
                <button
                  onClick={() => onSelectFriend(friend)}
                  className={`w-full p-3 text-left hover:bg-accent ${
                    selectedFriendId === friend.friend_id ? "bg-accent" : ""
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <span className={`inline-block h-2.5 w-2.5 shrink-0 rounded-full ${
                        onlineMap[friend.friend_id] ? "bg-green-500" : "bg-gray-400"
                      }`} />
                      <span className="font-medium">{friend.friend_name}</span>
                    </div>
                    {friend.unread_count > 0 && (
                      <span className="flex h-5 w-5 items-center justify-center rounded-full bg-primary text-xs text-primary-foreground">
                        {friend.unread_count}
                      </span>
                    )}
                  </div>
                  <div className="flex items-center justify-between text-sm text-muted-foreground">
                    <span className="truncate">{friend.last_message || "No messages"}</span>
                    <span>{formatTime(friend.last_message_at)}</span>
                  </div>
                </button>

                {openMenuId === friend.friend_id && !confirmUnfriendId && (
                  <div className="absolute right-2 top-2 z-10 w-40 rounded-md border bg-background p-1 shadow-lg">
                    <button
                      onClick={(e) => {
                        e.stopPropagation()
                        setConfirmUnfriendId(friend.friend_id)
                      }}
                      className="flex w-full items-center gap-2 rounded-sm px-2 py-1 text-sm text-destructive hover:bg-destructive/10"
                    >
                      <Trash2 className="h-4 w-4" />
                      Unfriend
                    </button>
                  </div>
                )}

                {confirmUnfriendId === friend.friend_id && (
                  <div className="absolute right-2 top-2 z-20 w-56 rounded-md border bg-background p-3 shadow-lg">
                    <div className="mb-2 text-sm font-medium">Unfriend {friend.friend_name}?</div>
                    <div className="flex gap-2">
                      <button
                        onClick={(e) => {
                          e.stopPropagation()
                          handleUnfriend(friend.friend_id, friend.friend_name)
                        }}
                        className="flex-1 rounded-sm bg-destructive px-2 py-1 text-sm text-destructive-foreground hover:bg-destructive/90"
                      >
                        Unfriend
                      </button>
                      <button
                        onClick={(e) => {
                          e.stopPropagation()
                          setConfirmUnfriendId(null)
                        }}
                        className="flex-1 rounded-sm border px-2 py-1 text-sm hover:bg-accent"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                )}

                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    setOpenMenuId(openMenuId === friend.friend_id ? null : friend.friend_id)
                  }}
                  className="absolute right-2 top-3 rounded-sm p-1 hover:bg-accent"
                >
                  <MoreVertical className="h-4 w-4" />
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}
