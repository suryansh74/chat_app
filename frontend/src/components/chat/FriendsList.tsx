import { useState, useEffect } from "react"
import { friendsApi, type FriendListItem } from "@/lib/api"

interface FriendsListProps {
  onSelectFriend: (friend: FriendListItem) => void
  selectedFriendId?: string
}

export function FriendsList({ onSelectFriend, selectedFriendId }: FriendsListProps) {
  const [friends, setFriends] = useState<FriendListItem[]>([])
  const [loading, setLoading] = useState(true)

  const loadFriends = async () => {
    setLoading(true)
    const { data } = await friendsApi.getFriends()
    if (data?.friends) {
      setFriends(data.friends)
    }
    setLoading(false)
  }

  useEffect(() => {
    loadFriends()
  }, [])

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
    <div className="flex h-full flex-col border-r">
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
              <li key={friend.friend_id}>
                <button
                  onClick={() => onSelectFriend(friend)}
                  className={`w-full p-3 text-left hover:bg-accent ${
                    selectedFriendId === friend.friend_id ? "bg-accent" : ""
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <span className="font-medium">{friend.friend_name}</span>
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
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}