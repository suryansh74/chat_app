import { useState, useCallback } from "react"
import { Navbar } from "@/components/Navbar"
import { FriendsList } from "@/components/chat/FriendsList"
import { ChatWindow } from "@/components/chat/ChatWindow"
import { FriendRequestModal } from "@/components/friends/FriendRequestModal"
import { NotificationPanel } from "@/components/notification/NotificationPanel"
import { useAuth } from "@/contexts/AuthContext"
import type { FriendListItem } from "@/lib/api"

export function HomePage() {
  const { user } = useAuth()
  const [selectedFriend, setSelectedFriend] = useState<FriendListItem | null>(null)
  const [showFriendRequest, setShowFriendRequest] = useState(false)
  const [showNotifications, setShowNotifications] = useState(false)
  const [friendsKey, setFriendsKey] = useState(0)

  const handleSelectFriend = (friend: FriendListItem) => {
    setSelectedFriend(friend)
  }

const handleNotificationsRefresh = useCallback(() => {
    setFriendsKey(prev => prev + 1)
  }, [])

  return (
    <div className="flex h-screen flex-col">
      <Navbar 
        onOpenFriendRequest={() => setShowFriendRequest(true)} 
        onOpenNotifications={() => setShowNotifications(true)}
        onNotificationCountChange={handleNotificationsRefresh}
        onNotificationHandled={handleNotificationsRefresh}
      />
      
      <main className="flex flex-1 overflow-hidden">
        {/* Left Sidebar - Friends List */}
        <div className="w-80 flex-shrink-0">
          <FriendsList 
            key={friendsKey}
            onSelectFriend={handleSelectFriend} 
            selectedFriendId={selectedFriend?.friend_id}
          />
        </div>

        {/* Right Side - Chat Window or Notifications */}
        <div className="flex-1">
          {showNotifications ? (
            <NotificationPanel 
              onClose={() => setShowNotifications(false)} 
              onNotificationHandled={handleNotificationsRefresh}
            />
          ) : selectedFriend && user ? (
            <ChatWindow friend={selectedFriend} userId={user.id} />
          ) : (
            <div className="flex h-full items-center justify-center text-muted-foreground">
              <div className="text-center">
                <p className="text-lg">Select a friend to start chatting</p>
                <p className="text-sm">Or add new friends using the + button in the navbar</p>
              </div>
            </div>
          )}
        </div>
      </main>

      {/* Friend Request Modal */}
      {showFriendRequest && (
        <FriendRequestModal onClose={() => setShowFriendRequest(false)} />
      )}
    </div>
  )
}