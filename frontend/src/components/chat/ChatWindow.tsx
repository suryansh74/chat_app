import { useState, useEffect, useRef, useCallback } from "react"
import { Send } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { chatApi, type MessageListItem, type FriendListItem } from "@/lib/api"
import { useWebSocket } from "@/contexts/WebSocketContext"

interface ChatWindowProps {
  friend: FriendListItem
  userId: string
}

export function ChatWindow({ friend, userId, onFriendRemoved }: ChatWindowProps & { onFriendRemoved: () => void }) {
  const { subscribe, unsubscribe } = useWebSocket()
  const [messages, setMessages] = useState<MessageListItem[]>([])
  const [newMessage, setNewMessage] = useState("")
  const [loading, setLoading] = useState(true)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const loadMessages = useCallback(async () => {
    setLoading(true)
    const { data } = await chatApi.getMessages(friend.friend_id, 50)
    if (data?.messages) {
      setMessages(data.messages.reverse())
    }
    setLoading(false)
  }, [friend.friend_id])

  useEffect(() => {
    loadMessages()
  }, [loadMessages])

  useEffect(() => {
    const handleFriendRemoved = (data: unknown) => {
      const d = data as { message: { user_id: string } }
      if (d.message.user_id === friend.friend_id) {
        onFriendRemoved()
      }
    }

    subscribe("friend_removed", handleFriendRemoved)
    return () => unsubscribe("friend_removed", handleFriendRemoved)
  }, [friend.friend_id, subscribe, unsubscribe, onFriendRemoved])

  useEffect(() => {
    const handleNewMessage = (data: unknown) => {
      const raw = data as { message: { id: string; from_user_id: string; to_user_id: string; content: string; created_at: string } }
      const msg = raw.message
      const isMe = msg.from_user_id === userId
      if (!isMe) {
        setMessages((prev) => [...prev, {
          message_id: msg.id,
          from_user_id: msg.from_user_id,
          to_user_id: msg.to_user_id,
          content: msg.content,
          is_me: false,
          created_at: msg.created_at,
        }])
      }
    }

    subscribe("new_message", handleNewMessage)
    return () => unsubscribe("new_message", handleNewMessage)
  }, [userId, subscribe, unsubscribe])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [messages])

  const sendMessage = async () => {
    if (!newMessage.trim()) return

    const { data } = await chatApi.sendMessage(friend.friend_id, newMessage)
    if (data?.message) {
      const raw = data.message as { id?: string; from_user_id?: string; to_user_id?: string; content?: string; created_at?: string; message_id?: string; is_me?: boolean }
      setMessages((prev) => [...prev, {
        message_id: raw.message_id || raw.id || "",
        from_user_id: raw.from_user_id || "",
        to_user_id: raw.to_user_id || "",
        content: raw.content || "",
        is_me: true,
        created_at: raw.created_at || new Date().toISOString(),
      }])
      setNewMessage("")
    }
  }

  const formatTime = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
  }

  return (
    <div className="flex h-full flex-col">
      <div className="border-b p-3">
        <h2 className="font-semibold">{friend.friend_name}</h2>
        <p className="text-sm text-muted-foreground">{friend.friend_email}</p>
      </div>

      <div className="flex-1 overflow-y-auto p-3">
        {loading ? (
          <div className="flex items-center justify-center h-full">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          </div>
        ) : messages.length === 0 ? (
          <div className="flex items-center justify-center h-full text-muted-foreground">
            No messages yet. Start the conversation!
          </div>
        ) : (
          <div className="space-y-2">
            {messages.map((msg) => (
              <div
                key={msg.message_id}
                className={`flex ${msg.is_me ? "justify-end" : "justify-start"}`}
              >
                <div
                  className={`max-w-[70%] rounded-lg p-2 ${
                    msg.is_me
                      ? "bg-primary text-primary-foreground"
                      : "bg-muted"
                  }`}
                >
                  <p>{msg.content}</p>
                  <p className="text-xs opacity-70">{formatTime(msg.created_at)}</p>
                </div>
              </div>
            ))}
            <div ref={messagesEndRef} />
          </div>
        )}
      </div>

      <div className="border-t p-3">
        <form
          onSubmit={(e) => {
            e.preventDefault()
            sendMessage()
          }}
          className="flex gap-2"
        >
          <Input
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Type a message..."
          />
          <Button type="submit" size="icon">
            <Send className="h-4 w-4" />
          </Button>
        </form>
      </div>
    </div>
  )
}
